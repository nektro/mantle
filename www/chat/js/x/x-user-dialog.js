"use strict";

//
import { create_element, dcTN } from "./../util.js";
// import { output } from "./../ui.util.js";
// import * as ui from "./../ui.js";
import * as api from "./../api/index.js";

//
customElements.define("x-user-dialog", class extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        //
    }
    /**
     * @param {string} uid
     * @param {MouseEvent} e
     */
    async openWith(uid, e) {
        const userN = await api.M.users.get(uid);
        this.removeAllChildren();
        this.appendChild(create_element("div", null, [
            create_element("h2", null, [dcTN(userN.nickname)]),
            create_element("h2", null, [
                create_element("span", null, [dcTN(userN.name)]),
                dcTN("#"),
                create_element("span", null, [dcTN(userN.id)]),
            ]),
            create_element("div", null, [dcTN(userN.uuid)]),
            create_element("div", null, [
                dcTN("Provider: "),
                create_element("span", null, [dcTN(userN.provider)]),
            ]),
            create_element("hr"),
            create_element("div", null, [dcTN("Roles")]),
            create_element("ol"),
            create_element("div", [["data-requires","manage_roles"]], [
                create_element("i", [["class","plus icon"]]),
                create_element("ol"),
            ]),
        ]));
        set_x(this.children[0], e.x);
        set_y(this.children[0], e.y);
    }
});

/**
 * @param {HTMLDivElement} el
 * @param {number} ex
 */
function set_x(el, ex) {
    let x = ex + 24;
    // const ew = el.offsetWidth;
    // const ww = window.inneroffsetWidth;
    // if (x + ew > ww) { x = ww - ew - 24; }
    el.style.left = `${x}px`;
}

/**
 * @param {HTMLDivElement} el
 * @param {number} ey
 */
function set_y(el, ey) {
    let y = ey - 24;
    const eh = el.offsetHeight;
    const wh = window.innerHeight;
    if (y + eh > wh) { y = wh - eh - 24; }
    el.style.top = `${y}px`;
}
