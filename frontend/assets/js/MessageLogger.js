export default class MessageLogger {
    constructor() {
        this.msgContainer = document.getElementById("msgContainer")
    }

    info(msg) {
        this.public('white', "info", msg)
    }

    error(msg) {
        this.public('red', 'error', msg)
    }

    public(color, preffix, msg) {
        let block = document.createElement('p')
        block.setAttribute('style', "color:" + color)
        block.innerHTML = `[${preffix}] ${msg}`

        let firstChild = this.msgContainer.firstChild
        this.msgContainer.insertBefore(block, firstChild)
    }
}