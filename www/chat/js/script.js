"use strict";
//
import "./x/index.js";
//
import { setDataBinding } from "./util.js";
import * as ui from "./ui.js";
import { el_2, el_3, el_1, output } from "./ui.util.js";
import * as api from "./api/index.js";
import * as ws from "./ws.js";

//
$(document).on("click", (e) => {
    const p = e.target.path();
    if (p.filter((v) => v.matches("dialog[open]")).length > 0) return;
    if (p.filter((v) => v.matches(".msg .usr")).length > 0) return;
    const s = document.querySelectorAll("dialog[open]");
    s.forEach((v) => v.removeAttribute("open"));
});

//
(async function() {
    //
    moment.defaultFormat = "ddd MMM DD Y HH:mm:ss zZZ";

    //
    await api.M.meta.about().then((x) => {
        setDataBinding("server_name", x.name);
        setDataBinding("server_version", x.version);
        //
        const sx = document.querySelector("x-settings[data-s-for=server] [data-s-section=overview]");
        for (const item of ["name","description","cover_photo","profile_photo","public"]) {
            sx.querySelector(`[name="${item}"]`).setAttribute("value", x[item]);
        }
        //
        el_2.children[1].addEventListener("click", () => {
            document.querySelector("x-settings[data-s-for=server]").setAttribute("open","");
        });
    });

    //
    await api.M.users.me().then((x) => {
        ui.volatile.me = x.user;
        ui.volatile.me.perms = x.perms;
        const n = ui.volatile.me.nickname || ui.volatile.me.name;
        el_3.children[0].textContent = `@${n}`;
        document.querySelectorAll("[data-requires]").forEach((el) => {
            el.setAttribute("hidden", "");
        });
        const p = x.perms;
        for (const key in p) {
            if (p[key]) {
                document.querySelectorAll(`[data-requires^="${key}"]`).forEach((el) => {
                    el.removeAttribute("hidden");
                });
            }
        }
        //
        el_3.children[1].addEventListener("click", () => {
            document.querySelector("x-settings[data-s-for=user]").setAttribute("open","");
        });
        //
        document.querySelectorAll("[data-s-for='user'] [data-s-section='my_account'] [fill]").forEach((el) => {
            el.setAttribute("fill", x.user.uuid);
        });
    }).catch(() => {
        location.assign("../");
    });

    //
    await api.M.roles.me().then((x) => {
        const rls = x.sort((a,b) => a.position > b.position);
        //
        for (const item of rls) {
            ui.M.role.add(item);
        }
    });

    //
    await api.M.channels.me().then(async (x) => {
        for (const item of x) {
            await ui.addChannel(item);
        }
        await output.setActiveChannel(x[0].uuid);

        const el2 = document.getElementById("channel-name");
        el2.children[0].textContent = x[0].name;
        el2.children[1].textContent = x[0].description;

        el_1.querySelector("button").addEventListener("click", async () => {
            const {value: name} = await Swal({
                title: "Enter the new channel's name",
                input: "text",
                showCancelButton: true,
                inputValidator: (value) => !value && "You need to write something!",
            });
            if (name !== undefined) {
                return api.M.channels.create(name);
            }
        });
        el_1.querySelector("ol").addEventListener("click", (ev) => {
            const fl = ev.composedPath().filter((v) => v instanceof Element && v.matches("[data-uuid]"));
            if (fl.length === 0) return;
            output.setActiveChannel(fl[0].dataset.uuid);
        });
    });

    await api.M.invites.me().then((x) => {
        for (const item of x) {
            ui.M.invite.add(item);
        }
    });

    await api.M.users.online().then((x) => {
        for (const item of x) {
            ui.M.user.connect(item);
        }
    });

    //
    /** @type {HTMLInputElement} */
    const input = document.getElementById("input").children[0];
    const socket = new WebSocket(`ws${location.protocol.substring(4)}//${location.host}/ws`);

    socket.addEventListener("open", () => {
        el_2.children[0].classList.remove("loading");
        el_2.children[0].classList.add("online");
    });
    socket.addEventListener("close", () => {
        el_2.children[0].classList.remove("online");
    });
    socket.addEventListener("message", async (e) => {
        const d = JSON.parse(e.data);
        let o = ws.M;
        for (const item of d.type.split("-")) {
            if (!(item in o)) {
                console.error("event handler not found:", d);
                return;
            }
            o = o[item];
        }
        if (typeof o !== "function") {
            console.error("handler is not a function:", `"${d.type}"`, `"${typeof o}"`);
            return;
        }
        await o(d);
    });
    setInterval(() => {
        if (el_2.children[0].classList.contains("online")) {
            socket.send(JSON.stringify({
                type: "ping",
            }));
        }
    }, 30*1000);

    //
    input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            const msg_con = e.target.value;
            if (msg_con.length === 0) return;
            socket.send(JSON.stringify({
                type: "message",
                in: output.active_channel_uid,
                message: msg_con.trim(),
            }));
            e.target.value = "";
        }
    });
})();
