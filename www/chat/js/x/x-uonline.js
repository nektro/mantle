"use strict";
//
import { create_element } from "./../util.js";
import * as api from "./../api/index.js";

//
customElements.define("x-uonline", class extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this.addRole({ uuid:"", name:"Online", position:9999 });
    }
    addRole(o) {
        const nel = create_element("x-uonline-role", [["uuid",o.uuid],["name",o.name],["position",o.position]]);
        //
        if (this.children.length === 0) {
            this.appendChild(nel);
        } else {
            for (let i = 0; i < this.children.length; i++) {
                const item = this.children[i];
                if (item._pos > o.position) {
                    this.insertBefore(nel, item);
                }
            }
        }
    }
    removeRole(o) {
        console.log(o,this);
    }
    async addUser(uid) {
        const o = await api.M.users.get(uid);
        const d = await o.getHightestDistinguishedRoleUID();
        this.querySelector(`x-uonline-role[uuid="${d}"]`).addUser(uid);
    }
    removeUser(uid) {
        this.querySelector(`x-uonline-user[uuid="${uid}"]`).removeMe();
    }
});
