"use strict";

//
import { create_element, dcTN, popup_set_x, popup_set_y } from "./../../../js/util.js";
import { el_xucm } from "./../../../js/ui.util.js";
import * as api from "./../../../js/api/index.js";

//
customElements.define("x-user-contextmenu", class extends HTMLElement {
    constructor() {
        super();
        this.triggers = [];
        this.triggers.push(this.localName);
    }
    /**
     * @param {string} uid
     * @param {MouseEvent} e
     */
    async openWith(uid, e) {
        e.preventDefault();
        const user = await api.M.users.get(uid);
        this.removeAllChildren();
        this.appendChild(create_element("ul", [], [
            create_element("li", [], [create_element("a", [["href",`./../~${user.uuid}`],["target","_hello"]], [dcTN("Profile")]),]),
        ]));
        popup_set_x(this, e.x);
        popup_set_y(this, e.y);
    }
});

//
document.addEventListener("click", (e) => {
    const p = e.target.path();
    for (const item of el_xucm.triggers) {
        for (const jtem of p) {
            if (jtem.matches(item)) {
                return;
            }
        }
    }
    el_xucm.removeAllChildren();
});
