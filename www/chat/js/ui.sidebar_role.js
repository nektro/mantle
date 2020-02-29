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
        if (this.count === 0 && x > 0) this.el.previousElementSibling.style.visibility = "visible";
        if (this.count > 0 && x === 0) this.el.previousElementSibling.style.visibility = "hidden";
        this.el.previousElementSibling.dataset.count = x.toString();
    }
}
