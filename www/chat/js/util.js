"use strict";
//
export const output = document.getElementById("messages").children[0];
export const userCache = new Map();
/** @type {Map<String,Array} */
export const messageCache = new Map();
export const el_1 = document.getElementById("channel-list");
export const el_2 = document.getElementById("server-name");
export const el_3 = document.getElementById("me");
export const el_4 = document.getElementById("users-online-list").children[1];

//

/**
 * @param {String} name
 * @param {String[][]} attrs
 * @param {Node[]} children
 * @returns {HTMLElement}
 */
export function create_element(name, attrs, children) {
    const ele = document.createElement(name);
    (attrs || []).forEach((v) => { ele.setAttribute(v[0], v[1]); });
    (children || []).forEach((v) => { ele.appendChild(v); });
    return ele;
}

/**
 * @param {String} string
 * @returns {Text}
 */
export function dcTN(string) {
    return document.createTextNode(string);
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

/**
 * @param {Number} x1
 * @param {Number} x2
 * @returns {Number[]}
 */
export function numsBetween(x1, x2) {
    if (x1 === x2) return [x1];
    const res = [];
    //
    if (x1 > x2) {
        for (let i = x2; i <= x1; i++) {
            res.push(i);
        }
    }
    if (x2 > x1) {
        for (let i = x1; i <= x2; i++) {
            res.push(i);
        }
    }
    return res;
}
