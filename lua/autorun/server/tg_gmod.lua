// in order this to work in case where both the gmod server instance and the telegram bot is running on the same local network, please ensure to run the gmod server with -allowlocalhttp parameter!

local HooksList = { // remove the hook from the list to disable the logging
	["PlayerSay"] = PlSayHook,
	["PlayerDeath"] = PlDeathHook,
	["player_disconnect"] = PlDisconnectHook,
	["player_connect"] = PlConnectHook,
}

ConcommandLogs = true // set to false if you want to disable console commands logging

// works only for commands added with concommand.Add, see https://gmo dwiki.com/concommand.Run
concommand._Run = concommand._Run or concommand.Run
function concommand.Run(ply, cmd, args, arg_str)
	if IsValid(ply) and ply:IsPlayer() and (cmd ~= nil) and ConcommandLogs then
	if arg_str == "" then
		arg_str = "(null)"
	else 
		arg_str = arg_str
	end
	HTTP( {
		method = "POST",
		parameters = {
			cmd = cmd,
			plr = ply:Nick(),
			args = arg_str,
			typ = "ConCmd",
		},
		url = "http://127.0.0.1/handlelog"
	} )
	end
	return concommand._Run(ply, cmd, args, arg_str)
end

function PlSayHook(ply, txt)
	HTTP( {
		method = "POST",
		parameters = {
			msg = txt,
			plr = ply:Nick(),
			typ = "ChatMsg",
		},
		url = "http://127.0.0.1/handlelog"
	} )
end

function PlDeathHook(victim, inflictor, attacker)
	if !victim:IsValid() then
		victim = "(null)"
	else 
		victim = victim:Nick() .. " (" .. victim:SteamID() .. ")"
	end

	if inflictor == nil then
		inflictor = "(null)"
	end

	if inflictor:GetClass() == nil then
		inflictor = inflictor:GetName()
	else
		inflictor = inflictor:GetClass()
	end

	if attacker:IsPlayer() and attacker:IsValid() then
		attacker = attacker:Nick() .. " (" .. attacker:SteamID() .. ")"
	else
		attacker = "world/prop/null"
		inflictor = "(null)"
	end

	HTTP( {
		method = "POST",
		parameters = {
			victim = victim,
			attacker = attacker,
			inflictor = inflictor,
			typ = "PlayerDeath",
		},
		url = "http://127.0.0.1/handlelog"
	} )	
end

function PlDisconnectHook(data)
	HTTP( {
		method = "POST",
		parameters = {
			plr = data.name .. " (" .. data.networkid .. ")",
			reason = data.reason,
			typ = "PlayerDisconnect"
		},
		url = "http://127.0.0.1/handlelog"
	} )	
end

function PlConnectHook(data)
		HTTP( {
		method = "POST",
		parameters = {
			plr = data.name .. " (" .. data.networkid .. ")",
			ipaddr = data.address,
			typ = "PlayerConnect"
		},
		url = "http://127.0.0.1/handlelog"
	} )	
end

for hk, fn in pairs(HooksList) do
    hook.Remove(hk, "CustomHook_" .. hk)
end

gameevent.Listen("player_disconnect")
gameevent.Listen("player_connect")

for hk, fn in pairs(HooksList) do
    local HookName = "CustomHook_" .. hk
    if hk == "player_disconnect" then
        hook.Add("player_disconnect", HookName, function(data) PlDisconnectHook(data) end)
        continue
    end
    if hk == "player_connect" then
        hook.Add("player_connect", HookName, function(data) PlConnectHook(data) end)
        continue
    end
    hook.Add(hk, HookName, fn)
end
