"use strict";
//
import * as ui from "./ui.js";
import * as api from "./api/index.js";
import { Channel } from "./ui.channel.js";
import { output, el_uonline } from "./ui.util.js";
import { setDataBinding } from "./util.js";

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
            if (["add_role","remove_role"].includes(d.key)) {
                document.querySelectorAll(`dialog.popup.user ol [data-role="${d.value}"]`).forEach((v) => {
                    v.classList.toggle("active");
                });
                el_uonline.checkUserForSwitch(o.uuid);
            }
        },
    },
    channel: {
        new: (d) => {
            ui.addChannel(new api.Channel(d.channel));
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
                x.href = y[0]+"="+(parseInt(y[1],10)+1).toString();
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
};
