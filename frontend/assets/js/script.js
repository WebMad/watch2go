import Video from "./Video.js";
let videoObject = null;

$(document).ready(() => {
    let video = document.getElementById('video');
    let video_url = $('#video_url');
    videoObject = new Video(video, video_url, 'ws://127.0.0.1:25566/ws');
    videoObject.init();
});