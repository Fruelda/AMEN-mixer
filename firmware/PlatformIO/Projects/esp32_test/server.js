const WebSocket = require("ws");


const server = new WebSocket.Server({

    host: "0.0.0.0",

    port: 8081,

    path: "/ws"

});



console.log(
    "AMEN WebSocket running ws://0.0.0.0:8081/ws"
);




server.on(
    "connection",
    (socket, request) => {


        console.log(
            "ESP32 CONNECTED",
            request.socket.remoteAddress
        );



        // Test server -> ESP32

        const welcome = JSON.stringify({

            type: "welcome",

            message: "websocket konek"

        });



        console.log(
            "SEND:",
            welcome
        );



        socket.send(
            welcome
        );





        // kirim test volume setelah 3 detik

        setTimeout(
            () => {

                const volume = JSON.stringify({

                    type: "test.volume",

                    channel: 1,

                    volume: 50

                });



                console.log(
                    "SEND:",
                    volume
                );



                socket.send(
                    volume
                );



            },

            3000

        );





        // menerima data ESP32

        socket.on(
            "message",
            (data) => {


                console.log(
                    "MESSAGE:",
                    data.toString()
                );


            }

        );





        socket.on(
            "close",
            () => {

                console.log(
                    "ESP32 DISCONNECTED"
                );

            }

        );


    }

);