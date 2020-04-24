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
            let addd = false;
            for (let i = 0; i < this.children.length; i++) {
                const item = this.children[i];
                if (!addd && item._pos > o.position) {
                    this.insertBefore(nel, item);
                    i++;
                    addd = true;
                }
                if (addd) {
                    item.check_for_switches();
                }
            }
        }
    }
    async removeRole(uid) {
        const e = this.querySelector(`x-uonline-role[uuid="${uid}"]`);
        for (const item of e.getAllUsers()) {
            e.removeUser(item);
            await this.addUser(item);
        }
        e.remove();
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
