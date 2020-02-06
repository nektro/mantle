/**
 */
//
export const output = document.getElementById("messages").children[0];
export const userCache = new Map();
export const messageCache = new Map();
export const el_1 = document.getElementById("channel-list");
export const el_2 = document.getElementById("server-name");
export const el_3 = document.getElementById("me");
export const el_4 = document.getElementById("users-online-list").children[1];

//

export function create_element(name, attrs, children) {
    var ele = document.createElement(name);
    (attrs || []).forEach(function(v) { ele.setAttribute(v[0], v[1]); });
    (children || []).forEach(function(v) { ele.appendChild(v); });
    return ele;
}

export function dcTN(string) {
    return document.createTextNode(string);
}

export async function getUserFromUUID(uuid) {
    if (userCache.has(uuid)) {
        return userCache.get(uuid);
    }
    else {
        const req = await fetch(`/api/users/${uuid}`);
        const res = await req.json();
        if (res.success) {
            userCache.set(uuid, res.message);
            return res.message;
        }
        else {
            return null;
        }
    }
}
