"use strict";
//
import { el_1 } from "./ui.util.js";

//
// jshint -W078
export class Channel {
    /**
     * @param {string} uid UUID
     * @returns {Channel}
     */
    constructor(uid) {
        this.el = el_1.querySelector(`[data-uuid="${uid}"]`);
    }
    /**
     * @returns {Number}
     */
    get unread() {
        return parseInt(this.el.dataset.unread, 10);
    }
    /**
     * @param {Number} x
     */
    set unread(x) {
        this.el.dataset.unread = x.toString();
        this.el.querySelector(".unred").textContent = x.toString();
    }
    set p_name(n) {
        this.el.children[0].textContent = n;
    }
}
