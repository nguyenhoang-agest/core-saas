import { NextResponse } from 'next/server';

const BACKEND_URL =
  process.env.BACKEND_INTERNAL_URL ||
  (process.env.NEXT_PUBLIC_API_URL && process.env.NEXT_PUBLIC_API_URL.startsWith('http')
    ? process.env.NEXT_PUBLIC_API_URL
    : 'http://localhost:8080');

export async function POST(request: Request) {
  try {
    const body = await request.json();

    const backendRes = await fetch(`${BACKEND_URL}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    const data = await backendRes.json();

    return NextResponse.json(data, {
      status: backendRes.status,
    });
  } catch (error: any) {
    return NextResponse.json(
      {
        code: 'E020001',
        message: 'Internal proxy error connecting to backend auth service',
        error: error.message,
      },
      { status: 500 }
    );
  }
}
