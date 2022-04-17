"use strict";
//
import * as ui from "./ui.js";
import * as api from "./api/index.js";
import { Channel } from "./ui.channel.js";
import { output, el_uonline, context, audio_buffer_size, vc_user_list } from "./ui.util.js";
import { create_element, dcTN, setDataBinding } from "./util.js";

//
export const M = {
    pong: () => {
        // do nothing, keep connection alive
    },
    user: {
        connect: (d) => {
            ui.M.user.connect(d.user);
        },
        disconnect: (d) => {
            ui.M.user.disconnect(d.user);
        },
        update: (d) => {
            const o = new api.User(d.user);
            if (["add_role", "remove_role"].includes(d.key)) {
                document.querySelectorAll(`dialog.popup.user ol [data-role="${d.value}"]`).forEach((v) => {
                    v.classList.toggle("active");
                });
                el_uonline.checkUserForSwitch(o.uuid);
                output.refreshUser(o.uuid);
            }
            if (d.key === "nickname") {
                el_uonline.checkUserForSwitch(o.uuid);
                output.refreshUser(o.uuid);
            }
        },
    },
    channel: {
        new: (d) => {
            ui.M.channel.add(new api.Channel(d.channel));
        },
        update: (d) => {
            const c = new api.Channel(d.channel);
            if (["name"].includes(d.key)) {
                new Channel(c.uuid).p_name = d.value;
            }
            if (output.active_channel_uid === c.uuid) {
                setDataBinding("channel_name", c.name);
                setDataBinding("channel_description", c.description);
            }
        },
        delete: (d) => {
            api.M.channels.remove(d.channel);
            ui.M.channel.remove(d.channel);
        },
    },
    role: {
        new: (d) => {
            ui.M.role.add(new api.Role(d.role));
        },
        update: (d) => {
            const o = new api.Role(d.role);
            if (["color"].includes(d.key)) {
                const x = document.getElementById("link-role-color");
                const y = x.href.split("=");
                x.href = y[0] + "=" + (parseInt(y[1], 10) + 1).toString();
            }
            if (d.key === "distinguish") {
                if (d.value === "0") {
                    el_uonline.removeRole(d.role.uuid);
                }
                if (d.value === "1") {
                    el_uonline.addRole(o);
                }
            }
        },
        delete: (d) => {
            api.M.roles.remove(d.role);
            ui.M.role.remove(d.role);
        },
    },
    message: {
        new: async (d) => {
            const u = await api.M.users.get(d.message.author);
            const m = new api.Message(d.message, d.in);
            await output.addMessage(d.in, u, m);
        },
        delete: (d) => {
            const ch = output.getChannel(d.channel);
            for (const item of d.affected) {
                api.M.channels.with(d.channel).messages.remove(item);
                ch.removeMessage(item);
            }
        },
    },
    invite: {
        new: (d) => {
            ui.M.invite.add(new api.Invite(d.invite));
        },
        update: (d) => {
            new api.Invite(d.invite);
        },
        delete: (d) => {
            api.M.invites.remove(d.invite);
            ui.M.invite.remove(d.invite);
        },
    },
    voice: {
        connect: async function (o) {
            const uid = o.userId;
            const user = await api.M.users.get(uid);
            const name = user.getName() + "#" + user.id;
            vc_user_list.appendChild(create_element("li", [["data-uuid", user.uuid]], [dcTN(name)]));
        },
        disconnect: function (o) {
            Array.from(vc_user_list.children).forEach((v) => {
                if (v.dataset.uuid === o.userId) {
                    v.remove();
                }
            });
        },
        data: async function (o) {
            if (o.from === ui.volatile.me.uuid) { return; }
            const a = context.createBuffer(1, audio_buffer_size, context.sampleRate);
            const b = a.getChannelData(0);
            for (let i = 0; i < a.length; i++) {
                b[i] = o.data[i];
            }
            const c = context.createBufferSource();
            c.buffer = a;
            c.connect(context.destination);
            c.start();
        },
    },
};
