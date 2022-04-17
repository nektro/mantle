"use strict";
//
// jshint -W003
import { create_element, dcTN } from "./util.js";
import { el_1, el_uonline, getSettingsSelection, output } from "./ui.util.js";
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

export function refresh_permissions() {
    for (const el of document.querySelectorAll("[data-requires]")) {
        el.dataset.hidden |= 0;
        el.dataset.hidden += 1;
    }
    for (const key in volatile.me.perms) {
        const has_perm = volatile.me.perms[key];
        if (has_perm) {
            for (const el of document.querySelectorAll(`[data-requires="${key}"]`)) {
                el.dataset.hidden -= 1;
            }
        }
    }
}

export function refresh_members() {
    for (const el of document.querySelectorAll("[data-mustbe-member]")) {
        el.dataset.hidden |= 0;
        el.dataset.hidden += 1;

        if (api.C.users.get(el.dataset.mustbeMember).is_member) {
            el.dataset.hidden -= 1;
        }
    }
}

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
            el_1.firstElementChild.appendChild(create_element("li", [["data-uuid", o.uuid], ["data-unread", "0"]], [
                create_element("div", [], [dcTN(o.name)]),
                // create_element("div", [["class","ments"]], [dcTN("0")]),
                create_element("div", [["class", "unred"]], [dcTN("0")]),
            ]));
            getSettingsSelection("server", "channels").addItem(o);
            await api.M.channels.with(o.uuid).messages.latest();
        },
        remove: (uid) => {
            el_1.firstElementChild.querySelector(`li[data-uuid="${uid}"]`).remove();
            getSettingsSelection("server", "channels").removeItem(uid);
            output.setActiveChannel(getSettingsSelection("server", "channels").items()[0]);
        },
    },
    role: {
        add: (o) => {
            getSettingsSelection("server", "roles").addItem(o);
            //
            if (o.distinguish) {
                document.querySelector("x-uonline").addRole(o);
            }
        },
        remove: (uid) => {
            getSettingsSelection("server", "roles").removeItem(uid);
            el_uonline.removeRole(uid);
        },
    },
    message: {
    },
    invite: {
        add: (o) => {
            getSettingsSelection("server", "invites").addItem(o);
        },
        remove: (uid) => {
            getSettingsSelection("server", "invites").removeItem(uid);
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
