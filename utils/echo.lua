if not args or #args == 0 then
    print("Usage: echo <text>")
    return
end

local text = table.concat(args, " ")

print(text)