"use strict";
//

//
/** @type {HTMLElement} */
export const output = document.getElementById("messages").children[0];
/** @type {Map<String,Array} */
export const messageCache = new Map();
export const el_1 = document.getElementById("channel-list");
export const el_2 = document.getElementById("server-name");
export const el_3 = document.getElementById("me");
export const el_4 = document.getElementById("users-online-list");

export const msg_processors = [
    ["/shrug", "¯\\_(ツ)_/¯"],
    ["/tableflip", "(╯°□°）╯︵ ┻━┻"],
    ["/unflip", "┬─┬ ノ( ゜-゜ノ)"],
];
