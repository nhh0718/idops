import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

// GET /api/ssh/keys — list all SSH keys in ~/.ssh/
export async function GET() {
  try {
    const { stdout } = await execAsync(`${CLI_PATH} ssh keys --json`);
    const keys = JSON.parse(stdout);
    return NextResponse.json({ keys });
  } catch {
    return NextResponse.json({ keys: [], unavailable: true });
  }
}

// DELETE /api/ssh/keys?name=xxx — delete a key pair
export async function DELETE(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const name = searchParams.get("name");
    if (!name) {
      return NextResponse.json({ error: "name is required" }, { status: 400 });
    }
    await execAsync(`${CLI_PATH} ssh keys delete ${name}`);
    return NextResponse.json({ success: true });
  } catch (error: unknown) {
    const stderr = (error as { stderr?: string })?.stderr || "";
    return NextResponse.json({ error: stderr || "Xóa key thất bại", success: false });
  }
}
