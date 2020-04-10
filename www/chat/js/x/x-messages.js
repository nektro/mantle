"use strict";
//
import { create_element } from "./../util.js";
import { Channel } from "./../ui.channel.js";
import * as ui from "./../ui.js";

//
customElements.define("x-messages", class extends HTMLElement {
    constructor() {
        super();
    }
    setActiveChannel(uid) {
        new Channel(uid).unread = 0;
        this.removeAllChildren();
        this.appendChild(create_element("x-msg-pane", [["uuid",uid]]));
    }
    get active_channel_uid() {
        return this.children[0]._uid;
    }
    /**
     * @param {string} ch_uid
     * @param {api.User} user
     * @param {api.Message} msg
     */
    async addMessage(ch_uid, user, msg) {
        const ch_sb = new Channel(ch_uid);
        if (localStorage.getItem("notifications_messages") === "1") {
            if (ch_uid !== this.active_channel_uid || !ui.volatile.windowActive) {
                const n = new Notification(`${user.getName()} (#${ch_sb.p_name})`, {
                    body: msg.body,
                    tag: ch_uid,
                });
                console.log(n);
            }
        }
        if (this.active_channel_uid !== ch_uid) {
            ch_sb.unread++;
            return;
        }
        await this.children[0].appendMessage(user, msg);
    }
    getChannel(uid) {
        return this.querySelector(`x-msg-pane[uuid="${uid}"]`);
    }
});
