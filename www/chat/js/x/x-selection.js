"use strict";
//
import { create_element, dcTN, deActivateChild } from "./../util.js";
import * as api from "./../api/index.js";

//
customElements.define("x-selection", class SSelection extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        $(this.children[0]).sortable({
            // jshint -W098
            stop: (ev,ue) => {
                const c = this.parentElement.getAttribute("data-s-section");
                const a = ue.item[0];
                const uid = a.dataset.uid;
                const pN = a.indexOfMe()+1;
                api.M[c].update(uid, "position", pN);
            },
        });
        const w = this.children[0].querySelector("a.new");
        if (w !== null) {
            w.addEventListener("click", async () => {
                const c = this.parentElement.getAttribute("data-s-section");
                const {value: name} = await Swal({
                    title: "Enter the new role's name",
                    input: "text",
                    showCancelButton: true,
                    inputValidator: (value) => !value && "You need to write something!",
                });
                if (name === undefined) return;
                return api.M[c].create(name);
            });
        }
    }
    addItem(item) {
        const rlist = this.children[0];
        const oLen = rlist.children.length;
        const nEl = create_element("a", [["data-uid",item.uuid]], [dcTN(item.name)]);
        nEl.addEventListener("click", (e) => {
            const et = e.target;
            this.setActive(Array.from(et.parentElement.children).indexOf(et));
        });
        rlist.insertBefore(
            nEl,
            rlist.querySelector(".div"),
        );
        if (oLen === 2) {
            rlist.parentElement.classList.add("active");
            this.setActive(0);
        }
    }
    setActive(i) {
        const rlist = this.children[0];
        deActivateChild(rlist);
        rlist.children[i].classList.add("active");
        const c = this.parentElement.getAttribute("data-s-section");
        const r = api.M[c].get(rlist.children[i].dataset.uid);
        const tin = rlist.parentElement.querySelectorAll("[fill]");
        for (const item of tin) {
            item.setAttribute("fill", r.uuid);
            const n = item.getAttribute("name");
            rlist.parentElement.querySelector(`[name="${n}"]`).setAttribute("value", r[n]);
        }
    }
});
