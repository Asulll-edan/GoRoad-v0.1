from typing import Optional
import asyncpg
from app.config import settings

pool: Optional[asyncpg.Pool] = None


async def init_db():
    global pool
    pool = await asyncpg.create_pool(
        host=settings.DB_HOST,
        port=settings.DB_PORT,
        user=settings.DB_USER,
        password=settings.DB_PASSWORD,
        database=settings.DB_NAME,
        min_size=2,
        max_size=10,
    )


async def close_db():
    global pool
    if pool:
        await pool.close()
        pool = None


async def fetch(query: str, *args):
    if not pool:
        return []
    async with pool.acquire() as conn:
        rows = await conn.fetch(query, *args)
        return [dict(r) for r in rows]


async def fetchrow(query: str, *args):
    if not pool:
        return None
    async with pool.acquire() as conn:
        row = await conn.fetchrow(query, *args)
        return dict(row) if row else None


async def execute(query: str, *args):
    if not pool:
        return
    async with pool.acquire() as conn:
        await conn.execute(query, *args)


async def executemany(query: str, args_list: list):
    if not pool:
        return
    async with pool.acquire() as conn:
        await conn.executemany(query, args_list)
