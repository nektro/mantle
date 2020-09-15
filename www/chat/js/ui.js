"use strict";
//
// jshint -W003
import { create_element, dcTN } from "./util.js";
import { el_1, el_uonline, output } from "./ui.util.js";
import * as api from "./api/index.js";

//
export const volatile = {
    /** @type {HTMLElement} */
    activeChannel: null,
    /** @type {Element[]} */
    selectedMsgs: [],
    /** @type {api.User} */
    me: null,
    windowActive: true,
};

//

export const M = {
    user: {
        connect: (uid) => {
            el_uonline.addUser(uid);
        },
        disconnect: (uid) => {
            el_uonline.removeUser(uid);
        },
    },
    channel: {
        add: async (o) => {
            el_1.firstElementChild.appendChild(create_element("li", [["data-uuid",o.uuid],["data-unread","0"]], [
                create_element("div", [], [dcTN(o.name)]),
                // create_element("div", [["class","ments"]], [dcTN("0")]),
                create_element("div", [["class","unred"]], [dcTN("0")]),
            ]));
            document.querySelector("x-settings[data-s-for=server] [data-s-section=channels] x-selection").addItem(o);
            await api.M.channels.with(o.uuid).messages.latest();
        },
        remove: (uid) => {
            el_1.firstElementChild.querySelector(`li[data-uuid="${uid}"]`).remove();
            document.querySelector("x-settings[data-s-for=server] [data-s-section=channels] x-selection").removeItem(uid);
            output.setActiveChannel(document.querySelector("x-settings[data-s-for=server] [data-s-section=channels] x-selection").items()[0]);
        },
    },
    role: {
        add: (o) => {
            document.querySelector("x-settings[data-s-for=server] [data-s-section=roles] x-selection").addItem(o);
            //
            if (o.distinguish) {
                document.querySelector("x-uonline").addRole(o);
            }
        },
        remove: (uid) => {
            document.querySelector("x-settings[data-s-for=server] [data-s-section=roles] x-selection").removeItem(uid);
            el_uonline.removeRole(uid);
        },
    },
    message: {
    },
    invite: {
        add: (o) => {
            document.querySelector("x-settings[data-s-for=server] [data-s-section=invites] x-selection").addItem(o);
        },
        remove: (uid) => {
            document.querySelector("x-settings[data-s-for=server] [data-s-section=invites] x-selection").removeItem(uid);
        },
    },
};

//

//
export const toggleHandlers = new Map();
//
function addToggleHandler(key_name, f) {
    toggleHandlers.set(key_name, (v) => {
        localStorage.setItem(key_name, v);
        f(v);
        document.querySelector(`x-settings [local-name="${key_name}"]`).setAttribute("value", v);
    });
}
addToggleHandler("notifications_messages", (v) => {
    if (v === "1") {
        Notification.requestPermission().then((result) => {
            if (result !== "granted") {
                toggleHandlers.get("notifications_messages")("0");
            }
        });
    }
});
