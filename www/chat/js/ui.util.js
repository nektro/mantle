"use strict";
//
import { numsNear } from "./util.js";

//
/** @type {HTMLElement} */
export const output = document.getElementById("messages").children[0];
export const userCache = new Map();
/** @type {Map<String,Array} */
export const messageCache = new Map();
export const el_1 = document.getElementById("channel-list");
export const el_2 = document.getElementById("server-name");
export const el_3 = document.getElementById("me");
export const el_4 = document.getElementById("users-online-list");

/**
 * @param {Element} ele an element.
 * @returns {Boolean} true if 'ele' is scrolled to within 5px of the bottom of its scroll.
 */
export function ele_atBottom(ele) {
    return numsNear(ele.scrollTop, ele.scrollHeight - ele.clientHeight, 5);
}

export async function getUserFromUUID(uuid) {
    if (userCache.has(uuid)) {
        return userCache.get(uuid);
    }
    const req = await fetch(`./../api/users/${uuid}`);
    const res = await req.json();
    if (!res.success) {
        return null;
    }
    userCache.set(uuid, res.message);
    return res.message;
}
