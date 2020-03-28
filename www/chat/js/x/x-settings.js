"use strict";
//
import { deActivateChild } from "./../util.js";

//
class SettingsDialog extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        if (this.dataset.active === undefined) {
            this.setActivePane(0);
        }
        for (const item of this.children[0].children[0].querySelectorAll("a:not(.div)")) {
            item.addEventListener("click", (ev) => {
                const t = ev.target;
                const i = Array.from(t.parentElement.querySelectorAll("a:not(.div)")).indexOf(t);
                this.setActivePane(i);
            });
        }
        this.addEventListener("click", (e) => {
            if (e.target.localName === "x-settings") {
                e.target.removeAttribute("open");
            }
        });
    nav() {
        return this.children[0].children[0];
    }
    }
    /**
     * @param {Number} n
     */
    setActivePane(n) {
        deActivateChild(this.children[0].children[1]);
        deActivateChild(this.nav());
        this.dataset.active = n.toString();
        this.children[0].children[0].querySelectorAll("a:not(.div)")[n].classList.add("active");
        this.children[0].children[1].children[n].classList.add("active");
    }
}

customElements.define("x-settings", SettingsDialog);
