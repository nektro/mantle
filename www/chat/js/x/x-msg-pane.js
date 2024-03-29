"use strict";
//
import { create_element, dcTN, ele_atBottom, safe_html_replace, setDataBinding } from "./../util.js";
import { async_ready, emoji, msg_processors } from "./../ui.util.js";
import * as api from "./../api/index.js";

//
/**
 * @param {api.User} user
 * @param {api.Message} msg
 */
async function _make_m_element(user, msg) {
    const attrsU = [["class", "usr"]];
    const rls = await user.getRoles();
    const a = rls.filter((v) => v.color.length > 0).sort((b, c) => b.position > c.position);
    if (a.length > 0) {
        attrsU.push(["data-role", a[0].uuid]);
    }
    //
    const tz = moment.tz.guess();
    const time = msg.time.tz(tz).format();
    const el = create_element("x-message", [["class", "msg"], ["uuid", msg.uuid], ["author", user.uuid]], [
        create_element("div", [["class", "ts"], ["title", time]], [dcTN(time.split(" ")[4])]),
        create_element("div", attrsU, [dcTN(user.getName())]),
        create_element("div", [["class", "dat"]], [dcTN(msg.body)]),
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
    safe_html_replace(mtx, /(`.+`)/gu, (match) => create_element("code", [], [dcTN(match.substring(1, match.length - 1))]));
    safe_html_replace(mtx, /(\|\|[\w ]+\|\|)/gu, (match) => create_element("x-spoiler", [], [dcTN(match.substring(2, match.length - 2))]));
    safe_html_replace(mtx, /([a-z]+:\/\/[^\s]+)/gu, (match) => create_element("a", [["href", match], ["target", "_blank"]], [dcTN(decodeURIComponent(match))]));
    safe_html_replace(mtx, /(magnet:[^\s]+)/gu, (match) => create_element("a", [["href", match], ["target", "_blank"]], [dcTN(decodeURIComponent(match))]));
    safe_html_replace(mtx, /(:[a-z0-9_+-]+:)/gu, (match) => {
        const n = match.substring(1, match.length - 1);
        const i = emoji.names.includes(n);
        return i ? dcTN(emoji.map[n]) : dcTN(match);
    });
    twemoji.parse(mtx);
    //
    return el;
}

//
/**
 * @param {api.Message} msg
 */
function _make_m_divider(msg) {
    return create_element("fieldset", [["class", "div date"]], [
        create_element("legend", [], [dcTN(msg.time.toString().substring(0, 15))])
    ]);
}

//
function _make_m_newdiv() {
    return create_element("fieldset", [["class", "div new"]], [
        create_element("legend", [], [dcTN("New")])
    ]);
}

//
customElements.define("x-msg-pane", class extends HTMLElement {
    constructor() {
        super();
    }

    async connectedCallback() {
        await async_ready;
        this._uid = this.getAttribute("uuid");
        //
        const c = await api.M.channels.get(this._uid);
        setDataBinding("channel_name", c.name);
        setDataBinding("channel_description", c.description);
        const hst = [...api.C.messages.get(this._uid)].map((v) => v[1]).sort((a, b) => a.id < b.id);
        for (const item of hst) {
            const u = await api.M.users.get(item.author);
            await this.prependMessage(u, item);
        }
        if (hst.length > 0 && hst.length < 50) {
            this.insertBefore(_make_m_divider(hst[hst.length - 1]), this.children[0]);
            this.classList.add("loading-done");
        }
        //
        this.addEventListener("scroll", async (e) => {
            if (this.children.length === 0) return;
            if (e.target.scrollTop !== 0) return;
            if (this.classList.contains("loading")) return;
            if (this.classList.contains("loading-done")) return;
            //
            this.classList.add("loading");
            const fc = this.firstMsgChild();
            await api.M.channels.with(this._uid).messages.after(fc._uid).then(async (y) => {
                if (y.length <= 1) {
                    this.classList.add("loading-done");
                    this.insertBefore(_make_m_divider(fc), this.firstElementChild);
                    return;
                }
                for (const item of y) {
                    const u = await api.M.users.get(item.author);
                    await this.prependMessage(u, item);
                }
                this.scrollTop = fc.offsetTop - 60;
            });
            this.classList.remove("loading");
        });
    }

    _scroll_to_bottom() {
        this.scrollTop = this.scrollHeight;
    }

    firstMsgChild() {
        for (const item of this.children) {
            if (item.classList.contains("msg")) {
                return item;
            }
        }
        return null;
    }

    /**
     * @param {api.User} user
     * @param {api.Message} msg
     */
    async appendMessage(user, msg, afk = false) {
        const at_bottom = ele_atBottom(this);
        //
        const prev_msg = this.lastElementChild;
        const ndv = this.querySelector(".div.new");
        //
        if (afk && ndv === null) {
            this.appendChild(_make_m_newdiv());
        }
        if (!afk && ndv !== null) {
            ndv.remove();
        }
        if (this.children.length === 0 || !prev_msg.time.isSame(msg.time, "day")) {
            this.appendChild(_make_m_divider(msg));
        }
        this.appendChild(await _make_m_element(user, msg));
        //
        if (at_bottom) this._scroll_to_bottom();
    }

    /**
     * @param {api.User} user
     * @param {api.Message} msg
     */
    async prependMessage(user, msg) {
        const at_bottom = ele_atBottom(this);
        //
        if (this.children.length === 0) {
            this.appendChild(await _make_m_element(user, msg));
            return;
        }
        const f = this.children[0];
        const d = f.classList.contains("date-div") ? this.children[1] : f;
        this.insertBefore(await _make_m_element(user, msg), f);
        if (!d.time.isSame(msg.time, "day")) {
            this.insertBefore(_make_m_divider(d), f);
        }
        //
        if (at_bottom) this._scroll_to_bottom();
    }

    selected() {
        return Array.from(this.querySelectorAll("x-message.selected"));
    }

    removeMessage(uid) {
        const el = this.querySelector(`x-message[uuid="${uid}"]`);
        el.remove();
    }

    async refreshUser(uid) {
        const u = await api.M.users.get(uid);
        const n = u.getName();
        const r = await u.getHightestColoredRoleUID();
        for (const item of this.children) {
            if (!item.classList.contains("msg")) { continue; }
            if (item._author !== uid) { continue; }
            const s = item.querySelector(".usr");
            s.textContent = n;
            s.dataset.role = r;
        }
    }
});
