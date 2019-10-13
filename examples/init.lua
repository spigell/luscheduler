local cron = require("cron")
local time = require("time")
local plugin = require("plugin")
local filepath = require("filepath")
local strings = require("strings")
local log = require("log")

local scenarios_dir = 'scenarios'

local scheduler = cron.new({verbose = 'true'})

for _, scenario_file in pairs(filepath.glob(scenarios_dir.."/*.lua")) do
  local file = io.open(scenario_file, "r")

  if (file) then
    for line in io.lines(scenario_file) do
      if strings.has_prefix(line, 'SCHEDULE') then
        get_schedule = loadstring(line)
        get_schedule()
      end
    end
  end
  if not (SCHEDULE == 'never') then
    print('[INFO] Set scenario ' .. scenario_file .. ' for ' .. SCHEDULE)
    scheduler:add_file(SCHEDULE, scenario_file)
  else
    print('[INFO] Skip scenario ' .. scenario_file)
  end
end

while true do
  time.sleep(600)
end

