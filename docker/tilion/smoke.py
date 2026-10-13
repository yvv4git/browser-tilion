import asyncio
import json
import sys

from tilion import Tilion


def _run(tilion, action, arg):
    coro_by_action = {
        "fetch": lambda: tilion.fetch(arg),
        "search": lambda: tilion.search(arg),
        "crawl": lambda: tilion.crawl(arg),
        "screenshot": lambda: tilion.screenshot(arg),
    }
    if action not in coro_by_action:
        raise SystemExit(f"unknown action: {action}")
    return coro_by_action[action]()


async def main():
    action = sys.argv[1]
    arg = sys.argv[2] if len(sys.argv) > 2 else ""
    async with Tilion(headless=True) as tilion:
        result = await _run(tilion, action, arg)
    print(json.dumps(result, indent=2, ensure_ascii=False, default=str))


asyncio.run(main())
