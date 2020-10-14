
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
            .coinslot {
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

            .coin {
                position: absolute;
                left: calc(100vw - 200px);
                top: 100px;
                height: 90px;
                color: #fff;
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
        <script>
            let paid = false;

            function onDrop(event) {
                var rect = document.getElementById('coinslot').getBoundingClientRect();
                console.log(rect.top, rect.right, rect.bottom, rect.left);

                console.log(event);
                console.log(event.clientX);
                console.log(event.clientY);

                if(event.clientX > (rect.left-100) && event.clientX < (rect.right+100) && event.clientY > (rect.top-100) && event.clientY < (rect.bottom+100)){
                    new Audio('./coin.mp3').play();
                    document.getElementById('coin').style.display = "none";
                    document.getElementById('words').style.display = "none";
                    document.getElementById('loading').style.display = "unset";
                    paid = true;
                }
            }

            function reject(){
                console.log('rejecting');
                if(paid){
                    document.getElementById('coin').style.display = "unset";
                    document.getElementById('words').style.display = "unset";
                    document.getElementById('loading').style.display = "none";
                }
            }
            
        </script>

        <div class="body">
            <div class="coinslot" style="z-index:1000" on:click={reject}>
                <img src="./images/coinslot.png" id="coinslot" width="100" height="100" />
            </div>
            <div class="coin" id="coin" draggable="true" on:dragend={onDrop}>
                <img src="./images/coin.png" class="rotate" width="100" height="100" />
            </div>
            <div class="words" id="words">
                <br/><br/><br/>Insert Coin
            </div>
            <div class="words" id="loading" style="display:none;">
                <br/><br/><br/>Loading
            </div>
        </div>


