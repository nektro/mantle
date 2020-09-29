"use strict";
//
import { create_element, dcTN } from "./../util.js";
import { el_uonline, el_xud } from "./../ui.util.js";
import * as api from "./../api/index.js";

//
customElements.define("x-uonline-user", class extends HTMLElement {
    constructor() {
        super();
    }
    async connectedCallback() {
        this._uid = this.getAttribute("uuid");
        const o = await api.M.users.get(this._uid);
        {
            const r = await o.getRoles();
            const l = r.filter((v) => v.color.length > 0);
            if (l.length > 0) {
                this.setAttribute("data-role", l[0].uuid);
            }
        }
        this.appendChild(create_element("span", null, [dcTN(o.getName())]));
        this.appendChild(create_element("span", null, [dcTN("#"+o.id)]));
        //
        el_xud.triggers.push("x-uonline-user");
        this.addEventListener("click", async (e) => {
            const target = e.target.path().filter((v) => v.tagName.toLowerCase() === "x-uonline-user")[0];
            el_xud.openWith(target._uid, e);
        });
    }
    get role_element() {
        return this.parentElement.parentElement;
    }
    removeMe() {
        this.role_element.removeUser(this._uid);
    }
    async check_for_switch() {
        const o = await api.M.users.get(this._uid);
        const r = await o.getHightestDistinguishedRoleUID();
        this.children[0].textContent = o.getName();
        if (r !== this.role_element._uid) {
            this.removeMe(o.uuid);
            el_uonline.addUser(o.uuid);
        }
    }
});
