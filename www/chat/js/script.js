"use strict";
//
import "./x/index.js";
//
import { create_element, dcTN, setDataBinding } from "./util.js";
import * as ui from "./ui.js";
import { el_2, el_3, el_1, output, messageCache, el_4 } from "./ui.util.js";
import * as api from "./api/index.js";

//
$("x-settings").on("click", (e) => {
    if (e.target.localName === "x-settings") {
        e.target.removeAttribute("open");
    }
});
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
        const n = ui.volatile.me.nickname || ui.volatile.me.name;
        el_3.children[0].textContent = `@${n}`;
        const p = x.perms;
        for (const key in p) {
            if (!p[key]) {
                document.querySelectorAll(`[data-requires^="${key}"]`).forEach((el) => {
                    el.setAttribute("hidden", "");
                });
            }
        }
        //
        el_3.children[1].addEventListener("click", () => {
            document.querySelector("x-settings[data-s-for=user]").setAttribute("open","");
        });
    }).catch(() => {
        location.assign("../");
    });

    //
    await api.M.roles.me().then((x) => {
        const rls = x.sort((a,b) => a.position > b.position);
        //
        for (const item of rls) {
            if (!item.distinguish) continue;
            el_4.appendChild(create_element("div", [["data-count","0"]], [dcTN(item.name)]));
            el_4.appendChild(create_element("ul", [["data-uid",item.uuid]], []));
        }
        el_4.appendChild(create_element("div", [["data-count","0"]], [dcTN("Online")]));
        el_4.appendChild(create_element("ul", [["data-uid",""]], []));
        //
        for (const item of rls) {
            ui.addRole(item);
        }
    });

    //
    await api.M.channels.me().then(async (x) => {
        for (const item of x) {
            await ui.addChannel(item);
        }
        await ui.setActiveChannel(x[0].uuid);

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
            ui.setActiveChannel(fl[0].dataset.uuid);
        });
        output.addEventListener("scroll", async (e) => {
            if (output.children.length === 0) return;
            if (e.target.scrollTop !== 0) return;
            if (output.classList.contains("loading")) return;
            if (output.classList.contains("loading-done")) return;
            //
            output.classList.add("loading");
            const fc = output.children[0];
            const lstm = output.children[0].dataset.msgUid;
            const chuid = ui.volatile.activeChannel.dataset.uuid;
            await fetch(`./../api/channels/${chuid}/messages?after=${lstm}`).then((y) => y.json()).then(async (y) => {
                if (y.message.length <= 1) {
                    output.classList.add("loading-done");
                    return;
                }
                for (let i = 1; i < y.message.length; i++) {
                    const item = y.message[i];
                    output.prepend(ui.createMessage(await api.M.users.get(item.author), item));
                    messageCache.get(chuid).unshift(item);
                }
                output.scrollTop = fc.offsetTop-60;
            });
            output.classList.remove("loading");
        });
        document.addEventListener("keydown", async (e) => {
            if (e.key !== "Delete") return;
            if (document.activeElement !== document.body) return;
            if (document.querySelector("dialog[open]") !== null) return;
            if (ui.volatile.selectedMsgs.length === 0) return;
            await Swal.fire({
                title: "Are you sure you want to delete?",
                text: "You won't be able to revert this!",
                type: "warning",
                showCancelButton: true,
            }).then(async (r) => {
                if (!r.value) return;
                const m2d = ui.volatile.selectedMsgs.filter((v) => v.dataset.userUid === ui.volatile.me.uuid).map((v) => v.dataset.msgUid);
                const fd = new FormData();
                m2d.forEach((v) => fd.append("ids", v));
                await fetch(`./../api/channels/${ui.volatile.activeChannel.dataset.uuid}/messages`, {
                    method: "DELETE",
                    body: fd,
                }).then((y) => y.json()).then(() => {
                    //
                });
            });
        });
    });

    await api.M.users.online().then((x) => {
        for (const item of x) {
            ui.setMemberOnline(item);
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
        switch (d.type) {
            case "pong": {
                // do nothing, keep connection alive
                break;
            }
            case "message": {
                const u = await api.M.users.get(d.message.author);
                ui.addMessage(d.in, u, d.message, true);
                break;
            }
            case "new-channel": {
                new api.Channel(d.channel);
                ui.addChannel(d.channel);
                break;
            }
            case "user-connect": {
                new api.User(d.user);
                ui.setMemberOnline(d.user);
                break;
            }
            case "user-disconnect": {
                ui.setMemberOffline(d.user);
                if (d.user === ui.volatile.me.uuid) socket.close();
                break;
            }
            case "message-delete": {
                const a1 = messageCache.get(d.channel);
                const a2 = a1.filter((v) => d.affected.includes(v.uuid));
                a2.forEach((v) => a1.splice(a1.indexOf(v), 1));
                if (output.dataset.active !== d.channel) break;
                for (const item of d.affected) {
                    const de = output.querySelector(`.msg[data-msg-uid="${item}"]`);
                    if (de !== null) de.remove();
                }
                break;
            }
            case "new-role": {
                new api.Role(d.role);
                ui.addRole(d.role);
                break;
            }
            case "user-update": {
                new api.User(d.user);
                if (["add_role","remove_role"].includes(d.key)) {
                    document.querySelectorAll(`dialog.popup.user ol [data-role="${d.value}"]`).forEach((v) => {
                        v.classList.toggle("active");
                    });
                }
                break;
            }
            case "role-update": {
                new api.Role(d.role);
                if (["color"].includes(d.key)) {
                    const x = document.getElementById("link-role-color");
                    const y = x.href.split("=");
                    x.href = y[0]+"="+(parseInt(y[1],10)+1).toString();
                }
                break;
            }
            default: {
                console.log(d);
            }
        }
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
                in: ui.volatile.activeChannel.dataset.uuid,
                message: msg_con.trim(),
            }));
            e.target.value = "";
        }
    });
})();
