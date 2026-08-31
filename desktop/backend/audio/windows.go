package audio

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/degubites/go-wca/pkg/wca"
	"github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

const hresultSFalse uintptr = 0x00000001

var (
	errWindowsAudioClosed = errors.New("windows audio worker is closed")

	iidAudioSessionManager2 = ole.NewGUID("{77AA99A0-1BD6-484F-8BC7-2C654C9A9B6F}")
	iidAudioSessionControl2 = ole.NewGUID("{BFB7FF88-7239-4FC9-8FA2-07C950BE9C6D}")
	iidSimpleAudioVolume    = ole.NewGUID("{87CE5498-68D6-44E5-9215-6DA47EF883D8}")
)

type windowsAudioRequest struct {
	run    func(*windowsAudioWorker) error
	result chan error
}

type windowsAudioWorker struct {
	enumerator *wca.IMMDeviceEnumerator
}

type WindowsAudio struct {
	requests chan windowsAudioRequest

	stop chan struct{}

	done chan struct{}

	closeOnce sync.Once

	initErr error
}

// =====================================================
// NEW WINDOWS AUDIO
// =====================================================

func NewWindowsAudio() *WindowsAudio {
	windowsAudio := &WindowsAudio{
		requests: make(chan windowsAudioRequest),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}

	ready := make(chan error, 1)

	go windowsAudio.run(ready)

	windowsAudio.initErr = <-ready

	if windowsAudio.initErr != nil {
		fmt.Printf(
			"[WINDOWS AUDIO] Core Audio init failed: %v\n",
			windowsAudio.initErr,
		)
	} else {
		fmt.Println(
			"[WINDOWS AUDIO] Core Audio ready",
		)
	}

	return windowsAudio
}

// =====================================================
// COM WORKER
// =====================================================

func (w *WindowsAudio) run(
	ready chan<- error,
) {
	runtime.LockOSThread()

	defer runtime.UnlockOSThread()
	defer close(w.done)

	if err := initializeCOM(); err != nil {
		ready <- fmt.Errorf(
			"CoInitializeEx(COINIT_MULTITHREADED): %w",
			err,
		)

		return
	}

	defer ole.CoUninitialize()

	var enumerator *wca.IMMDeviceEnumerator

	if err := wca.CoCreateInstance(
		wca.CLSID_MMDeviceEnumerator,
		0,
		wca.CLSCTX_ALL,
		wca.IID_IMMDeviceEnumerator,
		&enumerator,
	); err != nil {
		ready <- fmt.Errorf(
			"create IMMDeviceEnumerator: %w",
			err,
		)

		return
	}

	defer enumerator.Release()

	worker := &windowsAudioWorker{
		enumerator: enumerator,
	}

	ready <- nil

	for {
		select {
		case request := <-w.requests:
			request.result <- request.run(
				worker,
			)

		case <-w.stop:
			return
		}
	}
}

// =====================================================
// INITIALIZE COM
// =====================================================

func initializeCOM() error {
	err := ole.CoInitializeEx(
		0,
		ole.COINIT_MULTITHREADED,
	)

	if err == nil {
		return nil
	}

	if oleErr, ok := err.(*ole.OleError); ok &&
		oleErr.Code() == hresultSFalse {
		return nil
	}

	return err
}

// =====================================================
// EXECUTE ON COM THREAD
// =====================================================

func (w *WindowsAudio) do(
	run func(*windowsAudioWorker) error,
) error {
	if w == nil {
		return errWindowsAudioClosed
	}

	if w.initErr != nil {
		return w.initErr
	}

	request := windowsAudioRequest{
		run:    run,
		result: make(chan error, 1),
	}

	select {
	case w.requests <- request:

	case <-w.done:
		return errWindowsAudioClosed
	}

	select {
	case err := <-request.result:
		return err

	case <-w.done:
		return errWindowsAudioClosed
	}
}

// =====================================================
// DEFAULT RENDER DEVICE
// =====================================================

func (worker *windowsAudioWorker) withDefaultRenderDevice(
	run func(*wca.IMMDevice) error,
) error {
	var device *wca.IMMDevice

	if err := worker.enumerator.GetDefaultAudioEndpoint(
		wca.ERender,
		wca.EConsole,
		&device,
	); err != nil {
		return fmt.Errorf(
			"GetDefaultAudioEndpoint(eRender, eConsole): %w",
			err,
		)
	}

	defer device.Release()

	return run(
		device,
	)
}

// =====================================================
// MASTER ENDPOINT
// =====================================================

func (worker *windowsAudioWorker) withEndpointVolume(
	run func(*wca.IAudioEndpointVolume) error,
) error {
	return worker.withDefaultRenderDevice(
		func(device *wca.IMMDevice) error {
			var endpoint *wca.IAudioEndpointVolume

			if err := activateCOMInterface(
				device,
				wca.IID_IAudioEndpointVolume,
				unsafe.Pointer(&endpoint),
			); err != nil {
				return fmt.Errorf(
					"activate IAudioEndpointVolume: %w",
					err,
				)
			}

			defer endpoint.Release()

			return run(
				endpoint,
			)
		},
	)
}

// =====================================================
// AUDIO SESSION ENUMERATION
// =====================================================

func (worker *windowsAudioWorker) withMatchingSessions(
	app string,

	run func(
		pid uint32,
		processName string,
		volume *wca.ISimpleAudioVolume,
	) error,
) (
	int,
	[]string,
	error,
) {
	matched := 0

	seenSet := map[string]struct{}{}

	err := worker.withDefaultRenderDevice(
		func(device *wca.IMMDevice) error {
			// =========================================
			// SESSION MANAGER
			// =========================================

			var sessionManager *wca.IAudioSessionManager2

			if err := activateCOMInterface(
				device,
				iidAudioSessionManager2,
				unsafe.Pointer(&sessionManager),
			); err != nil {
				return fmt.Errorf(
					"activate IAudioSessionManager2: %w",
					err,
				)
			}

			defer sessionManager.Release()

			// =========================================
			// SESSION ENUMERATOR
			// =========================================

			var sessionEnumerator *wca.IAudioSessionEnumerator

			if err := sessionManager.GetSessionEnumerator(
				&sessionEnumerator,
			); err != nil {
				return fmt.Errorf(
					"IAudioSessionManager2.GetSessionEnumerator: %w",
					err,
				)
			}

			defer sessionEnumerator.Release()

			// =========================================
			// COUNT
			// =========================================

			count := 0

			if err := sessionEnumerator.GetCount(
				&count,
			); err != nil {
				return fmt.Errorf(
					"IAudioSessionEnumerator.GetCount: %w",
					err,
				)
			}

			// =========================================
			// ENUMERATE
			// =========================================

			for i := 0; i < count; i++ {
				if err := func() error {
					var sessionControl *wca.IAudioSessionControl

					if err := sessionEnumerator.GetSession(
						i,
						&sessionControl,
					); err != nil {
						return fmt.Errorf(
							"IAudioSessionEnumerator.GetSession(%d): %w",
							i,
							err,
						)
					}

					defer sessionControl.Release()

					// =================================
					// IAudioSessionControl2
					// =================================

					var sessionControl2 *wca.IAudioSessionControl2

					if err := sessionControl.PutQueryInterface(
						iidAudioSessionControl2,
						&sessionControl2,
					); err != nil {
						// System session atau session yang tidak
						// mendukung Control2 dilewati.
						return nil
					}

					defer sessionControl2.Release()

					// =================================
					// PID
					// =================================

					var pid uint32

					if err := sessionControl2.GetProcessId(
						&pid,
					); err != nil {
						return nil
					}

					if pid == 0 {
						return nil
					}

					// =================================
					// PROCESS NAME
					// =================================

					processName, err := processNameFromPID(
						pid,
					)

					if err != nil ||
						processName == "" {
						return nil
					}

					processName = strings.ToLower(
						processName,
					)

					seenSet[processName] = struct{}{}

					// =================================
					// MATCH APP
					// =================================

					if !appMatchesProcess(
						app,
						processName,
					) {
						return nil
					}

					// =================================
					// ISimpleAudioVolume
					// =================================

					var simpleVolume *wca.ISimpleAudioVolume

					if err := sessionControl.PutQueryInterface(
						iidSimpleAudioVolume,
						&simpleVolume,
					); err != nil {
						return fmt.Errorf(
							"QueryInterface(ISimpleAudioVolume) pid=%d process=%s: %w",
							pid,
							processName,
							err,
						)
					}

					defer simpleVolume.Release()

					matched++

					return run(
						pid,
						processName,
						simpleVolume,
					)
				}(); err != nil {
					return err
				}
			}

			return nil
		},
	)

	seen := make(
		[]string,
		0,
		len(seenSet),
	)

	for name := range seenSet {
		seen = append(
			seen,
			name,
		)
	}

	sort.Strings(
		seen,
	)

	return matched,
		seen,
		err
}

// =====================================================
// IMMDEVICE ACTIVATE
// =====================================================

func activateCOMInterface(
	device *wca.IMMDevice,
	iid *ole.GUID,
	target unsafe.Pointer,
) error {
	hr, _, _ := syscall.SyscallN(
		device.VTable().Activate,

		uintptr(
			unsafe.Pointer(
				device,
			),
		),

		uintptr(
			unsafe.Pointer(
				iid,
			),
		),

		uintptr(
			wca.CLSCTX_ALL,
		),

		0,

		uintptr(
			target,
		),
	)

	if hr != ole.S_OK {
		return ole.NewError(
			hr,
		)
	}

	return nil
}

// =====================================================
// REFRESH / SESSION DIAGNOSTIC
// =====================================================

func (w *WindowsAudio) Refresh() error {
	return w.do(
		func(worker *windowsAudioWorker) error {
			_, seen, err := worker.withMatchingSessions(
				"__amen_refresh_only__",

				func(
					uint32,
					string,
					*wca.ISimpleAudioVolume,
				) error {
					return nil
				},
			)

			if err != nil {
				return err
			}

			if len(seen) == 0 {
				fmt.Println(
					"[WINDOWS AUDIO] sessions: none",
				)

				return nil
			}

			fmt.Printf(
				"[WINDOWS AUDIO] sessions: %s\n",
				strings.Join(
					seen,
					", ",
				),
			)

			return nil
		},
	)
}

// =====================================================
// SET VOLUME
// =====================================================

func (w *WindowsAudio) SetVolume(
	app string,
	volume int,
) (
	int,
	error,
) {
	app = normalizeAppName(
		app,
	)

	volume = clampVolume(
		volume,
	)

	// =================================================
	// MASTER
	// =================================================

	if app == "master" {
		return w.setMasterVolume(
			volume,
		)
	}

	// =================================================
	// APPLICATION SESSION
	// =================================================

	actualVolume := volume

	matchedProcesses := map[string]struct{}{}

	matchedCount := 0

	var seen []string

	err := w.do(
		func(worker *windowsAudioWorker) error {
			var err error

			matchedCount,
				seen,
				err = worker.withMatchingSessions(
				app,

				func(
					pid uint32,
					processName string,
					simpleVolume *wca.ISimpleAudioVolume,
				) error {
					level := float32(
						volume,
					) / 100.0

					// =========================
					// SET SESSION VOLUME
					// =========================

					if err := simpleVolume.SetMasterVolume(
						level,
						nil,
					); err != nil {
						return fmt.Errorf(
							"ISimpleAudioVolume.SetMasterVolume app=%s pid=%d process=%s: %w",
							app,
							pid,
							processName,
							err,
						)
					}

					// =========================
					// READBACK
					// =========================

					var readback float32

					if err := simpleVolume.GetMasterVolume(
						&readback,
					); err != nil {
						return fmt.Errorf(
							"ISimpleAudioVolume.GetMasterVolume app=%s pid=%d process=%s: %w",
							app,
							pid,
							processName,
							err,
						)
					}

					actualVolume = scalarToPercent(
						readback,
					)

					matchedProcesses[processName] = struct{}{}

					return nil
				},
			)

			return err
		},
	)

	if err != nil {
		return volume,
			err
	}

	if matchedCount == 0 {
		return volume,
			noMatchingSessionError(
				app,
				seen,
			)
	}

	fmt.Printf(
		"[WINDOWS AUDIO] %s volume => %d%% sessions=%d processes=%s\n",
		app,
		actualVolume,
		matchedCount,
		strings.Join(
			sortedKeys(
				matchedProcesses,
			),
			",",
		),
	)

	return actualVolume,
		nil
}

// =====================================================
// MASTER VOLUME
// =====================================================

func (w *WindowsAudio) setMasterVolume(
	volume int,
) (
	int,
	error,
) {
	actualVolume := volume

	err := w.do(
		func(worker *windowsAudioWorker) error {
			return worker.withEndpointVolume(
				func(
					endpoint *wca.IAudioEndpointVolume,
				) error {
					level := float32(
						volume,
					) / 100.0

					if err := endpoint.SetMasterVolumeLevelScalar(
						level,
						nil,
					); err != nil {
						return fmt.Errorf(
							"SetMasterVolumeLevelScalar(%0.2f): %w",
							level,
							err,
						)
					}

					var readback float32

					if err := endpoint.GetMasterVolumeLevelScalar(
						&readback,
					); err != nil {
						return fmt.Errorf(
							"GetMasterVolumeLevelScalar: %w",
							err,
						)
					}

					actualVolume = scalarToPercent(
						readback,
					)

					return nil
				},
			)
		},
	)

	if err != nil {
		return volume,
			err
	}

	fmt.Printf(
		"[WINDOWS AUDIO] master volume => %d%%\n",
		actualVolume,
	)

	return actualVolume,
		nil
}

// =====================================================
// SET MUTE
// =====================================================

func (w *WindowsAudio) SetMute(
	app string,
	muted bool,
) (
	bool,
	error,
) {
	app = normalizeAppName(
		app,
	)

	// =================================================
	// MASTER
	// =================================================

	if app == "master" {
		return w.setMasterMute(
			muted,
		)
	}

	// =================================================
	// APPLICATION SESSION
	// =================================================

	actualMuted := muted

	matchedProcesses := map[string]struct{}{}

	matchedCount := 0

	var seen []string

	err := w.do(
		func(worker *windowsAudioWorker) error {
			var err error

			matchedCount,
				seen,
				err = worker.withMatchingSessions(
				app,

				func(
					pid uint32,
					processName string,
					simpleVolume *wca.ISimpleAudioVolume,
				) error {
					if err := simpleVolume.SetMute(
						muted,
						nil,
					); err != nil {
						return fmt.Errorf(
							"ISimpleAudioVolume.SetMute app=%s pid=%d process=%s: %w",
							app,
							pid,
							processName,
							err,
						)
					}

					if err := simpleVolume.GetMute(
						&actualMuted,
					); err != nil {
						return fmt.Errorf(
							"ISimpleAudioVolume.GetMute app=%s pid=%d process=%s: %w",
							app,
							pid,
							processName,
							err,
						)
					}

					matchedProcesses[processName] = struct{}{}

					return nil
				},
			)

			return err
		},
	)

	if err != nil {
		return muted,
			err
	}

	if matchedCount == 0 {
		return muted,
			noMatchingSessionError(
				app,
				seen,
			)
	}

	fmt.Printf(
		"[WINDOWS AUDIO] %s mute => %t sessions=%d processes=%s\n",
		app,
		actualMuted,
		matchedCount,
		strings.Join(
			sortedKeys(
				matchedProcesses,
			),
			",",
		),
	)

	return actualMuted,
		nil
}

// =====================================================
// MASTER MUTE
// =====================================================

func (w *WindowsAudio) setMasterMute(
	muted bool,
) (
	bool,
	error,
) {
	actualMuted := muted

	err := w.do(
		func(worker *windowsAudioWorker) error {
			return worker.withEndpointVolume(
				func(
					endpoint *wca.IAudioEndpointVolume,
				) error {
					if err := endpoint.SetMute(
						muted,
						nil,
					); err != nil {
						return fmt.Errorf(
							"SetMute(%t): %w",
							muted,
							err,
						)
					}

					if err := endpoint.GetMute(
						&actualMuted,
					); err != nil {
						return fmt.Errorf(
							"GetMute: %w",
							err,
						)
					}

					return nil
				},
			)
		},
	)

	if err != nil {
		return muted,
			err
	}

	fmt.Printf(
		"[WINDOWS AUDIO] master mute => %t\n",
		actualMuted,
	)

	return actualMuted,
		nil
}

// =====================================================
// GET MASTER STATE
// =====================================================

func (w *WindowsAudio) GetMasterState() (
	int,
	bool,
	error,
) {
	volume := 0

	muted := false

	err := w.do(
		func(worker *windowsAudioWorker) error {
			return worker.withEndpointVolume(
				func(
					endpoint *wca.IAudioEndpointVolume,
				) error {
					var level float32

					if err := endpoint.GetMasterVolumeLevelScalar(
						&level,
					); err != nil {
						return fmt.Errorf(
							"GetMasterVolumeLevelScalar: %w",
							err,
						)
					}

					if err := endpoint.GetMute(
						&muted,
					); err != nil {
						return fmt.Errorf(
							"GetMute: %w",
							err,
						)
					}

					volume = scalarToPercent(
						level,
					)

					return nil
				},
			)
		},
	)

	if err != nil {
		return 0,
			false,
			err
	}

	return volume,
		muted,
		nil
}

func (w *WindowsAudio) DebugSessions() error {
	return w.do(
		func(worker *windowsAudioWorker) error {

			_, seen, err :=
				worker.withMatchingSessions(
					"__debug__",

					func(
						pid uint32,
						processName string,
						volume *wca.ISimpleAudioVolume,
					) error {

						fmt.Printf(
							"[WINDOWS AUDIO SESSION] pid=%d process=%s\n",
							pid,
							processName,
						)

						return nil
					},
				)

			if err != nil {
				return err
			}

			fmt.Println(
				"[WINDOWS AUDIO SESSION LIST]",
				seen,
			)

			return nil
		},
	)
}

// =====================================================
// CLOSE
// =====================================================

func (w *WindowsAudio) Close() {
	if w == nil {
		return
	}

	w.closeOnce.Do(
		func() {
			close(
				w.stop,
			)
		},
	)

	<-w.done
}

// =====================================================
// NORMALIZE APP
// =====================================================

func normalizeAppName(
	app string,
) string {
	return strings.ToLower(
		strings.TrimSpace(
			app,
		),
	)
}

// =====================================================
// APP -> PROCESS MAPPING
// =====================================================

func appMatchesProcess(
	app string,
	processName string,
) bool {
	app = normalizeAppName(
		app,
	)

	processName = strings.ToLower(
		strings.TrimSpace(
			processName,
		),
	)

	switch app {
	case "browser":
		switch processName {
		case "chrome.exe",
			"msedge.exe",
			"firefox.exe",
			"brave.exe",
			"vivaldi.exe",
			"opera.exe",
			"browser.exe":

			return true

		default:
			return false
		}

	case "music":
		return processName ==
			"spotify.exe"

	case "chat":
		return processName ==
			"discord.exe" ||
			processName ==
				"discordcanary.exe" ||
			processName ==
				"discordptb.exe"
	case "game":
		return isGameProcess(processName)
	}

	if strings.HasSuffix(
		app,
		".exe",
	) {
		return processName ==
			app
	}

	return processName ==
		app ||
		processName ==
			app+".exe"
}

func isGameProcess(
	processName string,
) bool {
	knownGames := map[string]bool{
		"dota2.exe": true,
		"cs2.exe": true,
		"valorant-win64-shipping.exe": true,
		"eldenring.exe": true,
		"overwatch.exe": true,
	}

	return knownGames[processName]
}

// =====================================================
// PID -> PROCESS NAME
// =====================================================

func processNameFromPID(
	pid uint32,
) (
	string,
	error,
) {
	process, err := windows.OpenProcess(
		windows.PROCESS_QUERY_LIMITED_INFORMATION,
		false,
		pid,
	)

	if err != nil {
		return "",
			err
	}

	defer windows.CloseHandle(
		process,
	)

	buffer := make(
		[]uint16,
		32768,
	)

	size := uint32(
		len(buffer),
	)

	if err := windows.QueryFullProcessImageName(
		process,
		0,
		&buffer[0],
		&size,
	); err != nil {
		return "",
			err
	}

	fullPath := windows.UTF16ToString(
		buffer[:size],
	)

	return filepath.Base(
			fullPath,
		),
		nil
}

// =====================================================
// NO SESSION MATCH
// =====================================================

func noMatchingSessionError(
	app string,
	seen []string,
) error {
	if len(seen) == 0 {
		return fmt.Errorf(
			"no Windows audio sessions found on the default render endpoint for app %q",
			app,
		)
	}

	return fmt.Errorf(
		"no Windows audio session matched app %q; sessions seen: %s",
		app,
		strings.Join(
			seen,
			", ",
		),
	)
}

// =====================================================
// SORTED MAP KEYS
// =====================================================

func sortedKeys(
	values map[string]struct{},
) []string {
	result := make(
		[]string,
		0,
		len(values),
	)

	for value := range values {
		result = append(
			result,
			value,
		)
	}

	sort.Strings(
		result,
	)

	return result
}

// =====================================================
// CLAMP
// =====================================================

func clampVolume(
	volume int,
) int {
	if volume < 0 {
		return 0
	}

	if volume > 100 {
		return 100
	}

	return volume
}

// =====================================================
// SCALAR -> PERCENT
// =====================================================

func scalarToPercent(
	level float32,
) int {
	return clampVolume(
		int(
			math.Round(
				float64(level) *
					100.0,
			),
		),
	)
}
