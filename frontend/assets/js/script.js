import Video from "./Video.js";
let videoObject = null;

$(document).ready(() => {
    let video = document.getElementById('video');
    let video_url = $('#video_url');
    videoObject = new Video(video, video_url, 'ws://2.94.197.247:25565/ws');
    videoObject.init();
});