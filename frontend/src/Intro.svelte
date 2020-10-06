
        <link href='https://fonts.googleapis.com/css?family=Press+Start+2P' rel='stylesheet' type='text/css'>
        <style>
            .body {
                margin: 0;
                padding: 0;
                height: 100%;
                font-family: 'Press Start 2P', cursive;
                background: repeating-linear-gradient(to bottom, #0f0a1e, #0f0a1e 2px, #140e29 2px, #140e29 4px);
            }
            .words {
                position: absolute;
                top: 0;
                left: 0;
                bottom: 0;
                right: 0;
                height: 90px;
                color: #fff;
                width: 100%;
                text-align: center;
                margin: auto;
                font-size: 40px;
                line-height: 40px;
                letter-spacing: 5px;
                text-shadow: -2px 0 0 #fdff2a, -4px 0 0 #df4a42, 2px 0 0 #91fcfe, 4px 0 0 #4405fc;
                animation: blink 1s steps(4, start) infinite;
            }
            .coin {
                position: absolute;
                top: 0;
                left: 0;
                bottom: 0;
                right: 0;
                height: 90px;
                color: #fff;
                width: 100%;
                text-align: center;
                margin: auto;
            }
            @-moz-keyframes blink {
                to {
                    visibility: hidden;
                }
            }
            @-webkit-keyframes blink {
                to {
                    visibility: hidden;
                }
            }
            @-o-keyframes blink {
                to {
                    visibility: hidden;
                }
            }
            @keyframes blink {
                to {
                    visibility: hidden;
                }
            }

            .rotate {
                    -webkit-animation: rotate-vert-center 0.9s cubic-bezier(0.455, 0.030, 0.515, 0.955) infinite both;
                    animation: rotate-vert-center 0.9s cubic-bezier(0.455, 0.030, 0.515, 0.955) infinite both;
            }

            @-webkit-keyframes rotate-vert-center {
            0% {
                -webkit-transform: rotateY(0);
                        transform: rotateY(0);
            }
            100% {
                -webkit-transform: rotateY(360deg);
                        transform: rotateY(360deg);
            }
            }
            @keyframes rotate-vert-center {
                0% {
                    -webkit-transform: rotateY(0);
                            transform: rotateY(0);
                }
                100% {
                    -webkit-transform: rotateY(360deg);
                            transform: rotateY(360deg);
                }
            }
        </style>

        <div class="body">
            <div class="coin">
                <img src="./images/coin.png" class="rotate" width="100" height="100" />
            </div>
            <div class="words">
                <br/><br/><br/>Insert Coin
            </div>
        </div>


