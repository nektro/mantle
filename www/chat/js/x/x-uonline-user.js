"use strict";
//
import { create_element, dcTN } from "./../util.js";
import * as api from "./../api/index.js";

//
customElements.define("x-uonline-user", class extends HTMLElement {
    constructor() {
        super();
    }
    async connectedCallback() {
        this._uid = this.getAttribute("uuid");
        const o = await api.M.users.get(this._uid);
        this.appendChild(create_element("span", null, [dcTN(o.name)]));
        this.appendChild(create_element("span", null, [dcTN("#"+o.id)]));
    }
    removeMe() {
        this.parentElement.parentElement.removeUser(this._uid);
    }
});
