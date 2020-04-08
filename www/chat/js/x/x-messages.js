"use strict";
//
import { create_element } from "./../util.js";

//
customElements.define("x-messages", class extends HTMLElement {
    constructor() {
        super();
    }
    setActiveChannel(uid) {
        this.removeAllChildren();
        this.appendChild(create_element("x-msg-pane", [["uuid",uid]]));
    }
    get active_channel_uid() {
        return this.children[0]._uid;
    }
    async addMessage(ch_uid, user, msg) {
        if (this.active_channel_uid !== ch_uid) { return; }
        await this.children[0].appendMessage(user, msg);
    }
    getChannel(uid) {
        return this.querySelector(`x-msg-pane[uuid="${uid}"]`);
    }
});
