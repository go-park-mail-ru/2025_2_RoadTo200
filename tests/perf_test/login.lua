local thread_id_counter = 1
local threads = {}

function setup(thread)
    thread:set("thread_id", thread_id_counter)
    table.insert(threads, thread)
    thread_id_counter = thread_id_counter + 1
end

-- Функция инициализации скрипта
function init(args)
  request_counter = (thread_id - 1) * 10000000
end

-- Функция, вызываемая перед каждым запросом
request = function()
  local password = "qwertyui"

  -- Формируем тело запроса
  local body = string.format(
    '{"email": "user%d@test.ru", "password": "%s"}',
    request_counter + 2, password
  )

  -- Устанавливаем метод, заголовки и тело
  wrk.method = "POST"
  wrk.headers["Content-Type"] = "application/json"
  wrk.headers["User-Agent"] = "wrk-load-test"
  wrk.body = body

  return wrk.format()
end

-- Опционально: обработка ответов
response = function(status, headers, body)
  -- Можно логировать неудачные ответы
  if status ~= 200 then
    io.write(string.format("Error: %d\n", status))
  end
end
