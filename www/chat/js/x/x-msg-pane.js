"use strict";
//
import { create_element, dcTN, ele_atBottom } from "./../util.js";
import { msg_processors } from "./../ui.util.js";
import * as api from "./../api/index.js";

//
/**
 * @param {api.User} user
 * @param {api.Message} msg
 */
async function _make_m_element(user, msg) {
    const attrsU = [["class","usr"]];
    const rls = await user.getRoles();
    const a = rls.filter((v) => v.color.length > 0).sort((b,c) => b.position > c.position);
    if (a.length > 0) {
        attrsU.push(["data-role",a[0].uuid]);
    }
    //
    const tz = moment.tz.guess();
    const time = msg.time.tz(tz).format();
    const el = create_element("x-message", [["class","msg"],["uuid",msg.uuid],["author",user.uuid]], [
        create_element("div", [["class","ts"],["title",time]], [dcTN(time.split(" ")[4])]),
        create_element("div", attrsU, [dcTN(user.name)]),
        create_element("div", [["class","dat"]], [dcTN(msg.body)]),
    ]);
    const mtx = el.children[2];
    for (const item of msg_processors) {
        mtx.textContent = mtx.textContent.replace(item[0], () => {
            switch (typeof item[1]) {
                case "string": return item[1];
                case "function": return item[1]();
                default: break;
            }
            return "";
        });
    }
    mtx.innerHTML = mtx.textContent.replace(/(https?:\/\/[^\s]+)/gu, (match) => `<a target="_blank" href="${match}">${decodeURIComponent(match)}</a>`);
    twemoji.parse(mtx);
    //
    return el;
}

//
customElements.define("x-msg-pane", class extends HTMLElement {
    constructor() {
        super();
    }
    async connectedCallback() {
        this._uid = this.getAttribute("uuid");
        //
        for (const item of api.C.messages.get(this._uid)) {
            const u = await api.M.users.get(item[1].author);
            await this.prependMessage(u, item[1]);
        }
    }
    async appendMessage(user, msg) {
        const at_bottom = ele_atBottom(this);
        //
        this.appendChild(await _make_m_element(user, msg));
        //
        if (at_bottom) this.scrollTop = this.scrollHeight;
    }
    async prependMessage(user, msg) {
        const at_bottom = ele_atBottom(this);
        //
        if (this.children.length === 0) {
            this.appendChild(await _make_m_element(user, msg));
            return;
        }
        const f = this.children[0];
        this.insertBefore(await _make_m_element(user, msg), f);
        //
        if (at_bottom) this.scrollTop = this.scrollHeight;
    }
    selected() {
        return this.querySelectorAll("x-message.selected");
    }
});
