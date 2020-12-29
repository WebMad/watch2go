export default class Video {
    constructor(video_tag, src_input, ws_server) {
        this.video_tag = video_tag;
        this.src_input = $(src_input);
        this.ws_server = ws_server;
        this.src_input.on('change', () => {
            this.changeSrc(this.src_input.val());
        });

        this.pauseFromServer = false;

        this.src = "";
        this.time = null;
        this.is_paused = true;
        this.socket = null;
        this.keep_alive_interval = null;

        this.video_tag.onpause = () => this.setPause();
        this.video_tag.onplay = () => this.setPlay();
        this.video_tag.ontimeupdate = () => this.setTime();
    }

    init() {
        this.socket = new WebSocket(this.ws_server);
        this.socket.onmessage = (message) => this.onMessage(message);
        this.socket.onclose = () => {
            this.init();
        };
        this.socket.onopen = () => {
            if (this.keep_alive_interval) {
                clearInterval(this.keep_alive_interval);
            }
            this.socket.send(`{"type": "getData"}`);
            // this.keep_alive_interval = setInterval(() => this.socket.send(`{"type": "keepAlive"}`), 10000);
        };
    }

    onMessage(message) {
        let data = JSON.parse(message.data);
        this.time = data.time;
        if (data.isPaused !== this.video_tag.paused) {
            this.pauseFromServer = true
        }
        this.is_paused = data.isPaused;
        this.src = data.src;
    }

    setPause() {
        if (!this.pauseFromServer) {
            this.socket.send(`{"type": "pause","time": ${this.video_tag.currentTime}}`);
        }
        this.pauseFromServer = false
    }

    setPlay() {
        if (!this.pauseFromServer) {
            this.socket.send(`{"type": "play","time": ${this.video_tag.currentTime}}`);
        }
        this.pauseFromServer = false
    }

    setTime() {
        if (Math.abs(this._time - this.video_tag.currentTime) >= 2) {
            this.socket.send(`{"type": "changeTime","time": ${this.video_tag.currentTime}}`);
        }
        this._time = this.video_tag.currentTime;
    }

    changeSrc(src) {
        this.src = src;
        this.socket.send(`{"type": "src","src": "${src}"}`);
    }

    set src(value) {
        if (this._src !== value) {
            this.video_tag.src = value;
            this.src_input.val(value);
            if (this.video_tag.src.substr(-3, 3) == 'avi') {
                let source = document.createElement('source');
                source.src = value;
                source.type = "video/x-msvideo";
                video.appendChild(source)
            }
        }
        this._src = value;
    }

    set time(value) {
        this.video_tag.currentTime = value;
        this._time = value;
    }

    set is_paused(value) {
        if (value) {
            this.video_tag.pause();
        } else {
            this.video_tag.play();
        }
        this._is_paused = value;
    }
}