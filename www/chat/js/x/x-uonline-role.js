"use strict";
//
import { create_element, dcTN } from "./../util.js";

//
customElements.define("x-uonline-role", class extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this._uid = this.getAttribute("uuid");
        this._pos = this.getAttribute("position");
        this._name = this.getAttribute("name");
        this.appendChild(create_element("div", [["data-count","0"]], [dcTN(this._name)]));
        this.appendChild(create_element("ul"));
    }
    get count() {
        return parseInt(this.children[0].dataset.count, 10);
    }
    set count(x) {
        this.children[0].dataset.count = x.toString();
    }
    addUser(uid) {
        this.children[1].appendChild(create_element("x-uonline-user", [["uuid",uid]]));
        this.count += 1;
    }
    removeUser(uid) {
        this.querySelector(`x-uonline-user[uuid="${uid}"]`).remove();
        this.count -= 1;
    }
    getAllUsers() {
        return Array.from(this.children[1].children).map((v) => v._uid);
    }
    async check_for_switches() {
        for (const item of this.children[1].children) {
            await item.check_for_switch();
        }
    }
});
