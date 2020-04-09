"use strict";
//
// jshint -W003
import { create_element, dcTN } from "./util.js";
import { el_1, el_uonline } from "./ui.util.js";
import * as api from "./api/index.js";

//
export const volatile = {
    /** @type {HTMLElement} */
    activeChannel: null,
    /** @type {Element[]} */
    selectedMsgs: [],
    /** @type {api.User} */
    me: null,
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
    },
    role: {
        add: (o) => {
            document.querySelector("x-settings[data-s-for=server] [data-s-section=roles] x-selection").addItem(o);
            //
            if (o.distinguish) {
                document.querySelector("x-uonline").addRole(o);
            }
            //
            const nEl2 = create_element("li", [["data-role",o.uuid],["class","bg-bf"]], [dcTN(o.name)]);
            nEl2.addEventListener("click", (e) => {
                if (!volatile.me.perms.manage_roles) return;
                const et = e.target;
                const rid = et.dataset.role;
                const uid = document.querySelector("[data-bind=pp_user_uuid]").textContent;
                return api.M.users.update(uid,"remove_role",rid);
            });
            document.querySelector("dialog.popup.user ol").appendChild(nEl2);
            //
            const nEl3 = create_element("li", [["data-role",o.uuid]], [dcTN(o.name)]);
            nEl3.addEventListener("click", (e) => {
                const et = e.target;
                const rid = et.dataset.role;
                const uid = document.querySelector("[data-bind=pp_user_uuid]").textContent;
                return api.M.users.update(uid,"add_role",rid);
            });
            document.querySelector("dialog.popup.user div ol").appendChild(nEl3);
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

/**
 * @param {api.Channel} ch
 */
export async function addChannel(ch) {
    el_1.firstElementChild.appendChild(create_element("li", [["data-uuid",ch.uuid],["data-unread","0"]], [
        create_element("div", [], [dcTN(ch.name)]),
        create_element("div", [["class","unred"]], [dcTN("0")]),
    ]));
    //
    document.querySelector("x-settings[data-s-for=server] [data-s-section=channels] x-selection").addItem(ch);
    //
    await api.M.channels.with(ch.uuid).messages.latest();
}
