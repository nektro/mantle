"use strict";

//
import { create_element, dcTN, popup_set_x, popup_set_y } from "./../util.js";
import * as ui from "./../ui.js";
import * as api from "./../api/index.js";
import { el_xud } from "../ui.util.js";

//
customElements.define("x-user-dialog", class extends HTMLElement {
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
        const userN = await api.M.users.get(uid);
        this.removeAllChildren();
        this.appendChild(create_element("div", null, [
            create_element("h2", null, [dcTN(userN.nickname)]),
            create_element("h2", null, [
                create_element("span", null, [dcTN(userN.name)]),
                dcTN("#"),
                create_element("span", null, [dcTN(userN.id)]),
            ]),
            create_element("div", [["id", "pp_uuid"]], [dcTN(userN.uuid)]),
            create_element("div", null, [
                dcTN("Provider: "),
                create_element("span", null, [dcTN(userN.provider)]),
            ]),
            create_element("hr"),
            create_element("div", null, [dcTN("Roles")]),
            create_element("ol"),
        ]));
        for (const item of api.C.roles.values()) {
            if (item.id === undefined) { continue; }
            //
            const nEl2 = create_element("li", [["data-role", item.uuid], ["class", "bg-bf"]], [dcTN(item.name)]);
            nEl2.addEventListener("click", (ev) => {
                if (!ui.volatile.me.perms.manage_roles) return;
                const et = ev.target;
                const rid = et.dataset.role;
                const uid_role = document.querySelector("#pp_uuid").textContent;
                this.toggleRole(rid);
                return api.M.users.update(uid_role, "remove_role", rid);
            });
            this.children[0].querySelectorAll("ol")[0].appendChild(nEl2);
        }
        for (const item of await userN.getRoles()) {
            if (item.id === undefined) { continue; }
            this.toggleRole(item.uuid);
        }
        popup_set_x(this.children[0], e.x);
        popup_set_y(this.children[0], e.y);
    }

    toggleRole(uid) {
        for (const item of this.querySelectorAll("ol")) {
            item.querySelector(`[data-role="${uid}"]`).classList.toggle("active");
        }
    }
});

//
document.addEventListener("click", (e) => {
    const p = e.target.path();
    for (const item of el_xud.triggers) {
        for (const jtem of p) {
            if (jtem.matches(item)) {
                return;
            }
        }
    }
    el_xud.removeAllChildren();
});
