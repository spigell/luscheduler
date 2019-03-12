filepath = require('filepath')
local scenarios_dir = "examples/".."scenarios"

for _, scenario in pairs(filepath.glob(scenarios_dir.."/*.lua")) do

	local contents = ""
    local file = io.open( scenario, "r" )
    if (file) then
        contents = file:read()
        file:close()
        get_schedule = loadstring(contents)
        get_schedule()
    end
	schedule.new(SCHEDULE, scenario)
end

function sleep()
  os.execute("while sleep 3600; do :; done")
end
sleep()
