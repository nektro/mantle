"use strict";
//

export class SidebarRole {
    /**
     * @param {HTMLElement} el
     * @returns {SidebarRole}
     */
    constructor(el) {
        this.el = el;
    }
    /**
     * @returns {Number}
     */
    get count() {
        return parseInt(this.el.previousElementSibling.dataset.count, 10);
    }
    /**
     * @param {Number} x
     */
    set count(x) {
        this.el.previousElementSibling.dataset.count = x.toString();
    }
}
