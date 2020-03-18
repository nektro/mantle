"use strict";
//
import * as api from "./api/index.js";

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
export const roleCache = new Map();

roleCache.set("o", new api.Role({name:"Owner",color:""}));

export async function getUserFromUUID(uuid) {
    if (userCache.has(uuid)) {
        return userCache.get(uuid);
    }
    const u = await api.M.users.get(uuid);
    if (u.is_null) {
        return null;
    }
    userCache.set(uuid, u);
    return u;
}
