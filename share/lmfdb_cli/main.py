"""
LMFDB CLI - 命令行工具
"""
from __future__ import annotations

import json
from typing import Any, Annotated, Optional

import rich
import typer
from rich import print
from rich.console import Console
from rich.table import Table
from rich.tree import Tree

app = typer.Typer(
    name="lmfdb",
    help="LMFDB CLI - 查询 L-Functions and Modular Forms Database",
    add_completion=False,
)
console = Console()

# 常用 API 集合
API_COLLECTIONS = {
    "nf_fields": "Number fields (数域)",
    "ec_curvedata": "Elliptic curves (椭圆曲线)",
    "ec_classdata": "Elliptic curve isogeny classes",
    "g2c_curves": "Genus 2 curves (二 genus 曲线)",
    "char_dirichlet": "Dirichlet characters",
    "modlgal_reps": "Modular Galois representations",
    "maass_newforms": "Maass forms",
    "mf_newforms": "Modular forms",
    "lf_fields": "Local fields",
    "artin": "Artin representations",
}


@app.command()
def list_collections():
    """列出所有可用的 API 集合"""
    print("\n[bold]Available API Collections:[/bold]\n")
    
    tree = Tree("[bold]LMFDB Collections[/bold]")
    
    for name, desc in API_COLLECTIONS.items():
        tree.add(f"[cyan]{name}[/cyan] - {desc}")
    
    print(tree)
    print()


@app.command()
def nf(
    degree: Annotated[int, typer.Option("-d", "--degree", help="Number field degree")] = 2,
    discriminant: Annotated[Optional[int], typer.Option("--disc", help="Discriminant")] = None,
    class_number: Annotated[Optional[int], typer.Option("-h", "--class-number", help="Class number")] = None,
    limit: Annotated[int, typer.Option("-n", "--limit", help="Results limit")] = 10,
    fields: Annotated[Optional[str], typer.Option("-f", "--fields", help="Fields to return (comma-separated)")] = None,
    output: Annotated[Optional[str], typer.Option("-o", "--output", help="Output file (JSON)")] = None,
    headless: Annotated[bool, typer.Option("--headless/--no-headless", help="Run browser in headless mode")] = True,
):
    """
    查询 Number Fields (数域)
    
    Examples:
    
        lmfdb nf -d 2 -n 10
        lmfdb nf --discriminant -5
        lmfdb nf -d 3 --fields label,degree,disc
    """
    from lmfdb_cli.client import get_client, close_client
    
    try:
        client = get_client(headless=headless)
        
        # 构建查询参数
        params = {}
        if degree:
            params["degree"] = f"i{degree}"
        if discriminant:
            params["disc"] = f"i{discriminant}"
        if class_number:
            params["class_number"] = f"i{class_number}"
        
        # 解析字段
        field_list = None
        if fields:
            field_list = [f.strip() for f in fields.split(",")]
        
        # 执行查询
        console.print(f"[yellow]Querying LMFDB API...[/yellow]")
        
        data = client.query(
            collection="nf_fields",
            params=params,
            fields=field_list,
            limit=limit,
        )
        
        # 输出结果
        if output:
            with open(output, "w") as f:
                json.dump(data, f, indent=2)
            console.print(f"[green]Results saved to {output}[/green]")
        else:
            _display_nf_results(data, limit)
            
    finally:
        close_client()


@app.command()
def ec(
    rank: Annotated[Optional[int], typer.Option("-r", "--rank", help="Mordell-Weil rank")] = None,
    torsion: Annotated[Optional[int], typer.Option("-t", "--torsion", help="Torsion order")] = None,
    conductor: Annotated[Optional[int], typer.Option("--conductor", help="Conductor")] = None,
    limit: Annotated[int, typer.Option("-n", "--limit", help="Results limit")] = 10,
    fields: Annotated[Optional[str], typer.Option("-f", "--fields", help="Fields to return (comma-separated)")] = None,
    output: Annotated[Optional[str], typer.Option("-o", "--output", help="Output file (JSON)")] = None,
    headless: Annotated[bool, typer.Option("--headless/--no-headless", help="Run browser in headless mode")] = True,
):
    """
    查询 Elliptic Curves (椭圆曲线)
    
    Examples:
    
        lmfdb ec -r 2 -n 10
        lmfdb ec -t 5
        lmfdb ec --rank 0 --fields label,conductor,rank,torsion
    """
    from lmfdb_cli.client import get_client, close_client
    
    try:
        client = get_client(headless=headless)
        
        # 构建查询参数
        params = {}
        if rank is not None:
            params["rank"] = f"i{rank}"
        if torsion is not None:
            params["torsion"] = f"i{torsion}"
        if conductor:
            params["conductor"] = f"i{conductor}"
        
        # 解析字段
        field_list = None
        if fields:
            field_list = [f.strip() for f in fields.split(",")]
        
        # 执行查询
        console.print(f"[yellow]Querying LMFDB API...[/yellow]")
        
        data = client.query(
            collection="ec_curvedata",
            params=params,
            fields=field_list,
            limit=limit,
        )
        
        # 输出结果
        if output:
            with open(output, "w") as f:
                json.dump(data, f, indent=2)
            console.print(f"[green]Results saved to {output}[/green]")
        else:
            _display_ec_results(data, limit)
            
    finally:
        close_client()


@app.command()
def query(
    collection: Annotated[str, typer.Argument(help="API collection name")],
    limit: Annotated[int, typer.Option("-n", "--limit", help="Results limit")] = 10,
    fields: Annotated[Optional[str], typer.Option("-f", "--fields", help="Fields to return (comma-separated)")] = None,
    output: Annotated[Optional[str], typer.Option("-o", "--output", help="Output file (JSON)")] = None,
    headless: Annotated[bool, typer.Option("--headless/--no-headless", help="Run browser in headless mode")] = True,
    # Note: Additional kwargs are passed as query parameters
    kwargs_json: Annotated[Optional[str], typer.Option("-k", "--kwargs", help="Additional query params as JSON")] = None,
):
    """
    通用查询命令
    
    Examples:
    
        lmfdb query nf_fields -n 5
        lmfdb query ec_curvedata -k '{"rank": "i2"}' -n 20
    """
    from lmfdb_cli.client import get_client, close_client
    
    try:
        client = get_client(headless=headless)
        
        # 解析字段
        field_list = None
        if fields:
            field_list = [f.strip() for f in fields.split(",")]
        
        # 解析额外的查询参数
        params = {}
        if kwargs_json:
            import json as json_module
            params = json_module.loads(kwargs_json)
        
        # 执行查询
        console.print(f"[yellow]Querying {collection}...[/yellow]")
        
        data = client.query(
            collection=collection,
            params=params,
            fields=field_list,
            limit=limit,
        )
        
        # 输出结果
        if output:
            with open(output, "w") as f:
                json.dump(data, f, indent=2)
            console.print(f"[green]Results saved to {output}[/green]")
        else:
            _display_generic_results(data, collection, limit)
            
    finally:
        close_client()


def _display_nf_results(data: dict, limit: int):
    """显示 Number Fields 结果"""
    # 兼容 data 或 results 键
    items = data.get("data", data.get("results", []))[:limit]
    
    if not items:
        console.print("[yellow]No results found[/yellow]")
        return
    
    table = Table(title=f"Number Fields (showing {len(items)} of {len(items)} results)")
    
    # 动态添加列 - 使用第一个item的键
    if items and isinstance(items[0], dict):
        cols = list(items[0].keys())[:6]  # 最多6列
        
        for col in cols:
            table.add_column(col, style="cyan")
        
        for item in items:
            row = [str(item.get(col, "N/A"))[:20] for col in cols]  # 截断长字段
            table.add_row(*row)
    else:
        for item in items:
            console.print(item)
    
    console.print(table)


def _display_ec_results(data: dict, limit: int):
    """显示 Elliptic Curves 结果"""
    items = data.get("data", data.get("results", []))[:limit]
    
    if not items:
        console.print("[yellow]No results found[/yellow]")
        return
    
    if items and isinstance(items[0], dict):
        table = Table(title=f"Elliptic Curves (showing {len(items)} results)")
        
        cols = list(items[0].keys())[:6]
        
        for col in cols:
            table.add_column(col, style="cyan")
        
        for item in items:
            row = [str(item.get(col, "N/A"))[:20] for col in cols]
            table.add_row(*row)
        
        console.print(table)
    else:
        for item in items:
            console.print(item)


def _display_generic_results(data: dict, collection: str, limit: int):
    """显示通用结果"""
    # 检查是否是错误响应
    if "raw" in data and "Unable to locate" in str(data.get("raw", "")):
        console.print("[red]Error:[/red] Invalid collection or field name")
        console.print(f"URL: {data.get('url', 'N/A')}")
        return
    
    if "data" in data:
        items = data.get("data", [])[:limit]
        
        if not items:
            console.print("[yellow]No results found[/yellow]")
            return
        
        # 尝试显示表格
        if items and isinstance(items[0], dict):
            table = Table(title=f"{collection} Results")
            
            # 添加所有字段列
            for key in items[0].keys():
                table.add_column(key, style="cyan")
            
            for item in items:
                row = [str(item.get(key, "N/A"))[:20] for key in items[0].keys()]
                table.add_row(*row)
            
            console.print(table)
        else:
            for item in items:
                console.print(item)
    else:
        # 可能是错误或原始 HTML
        console.print(data)


@app.command()
def version():
    """显示版本信息"""
    console.print("[bold]LMFDB CLI[/bold] version 0.1.0")
    console.print("Query LMFDB from command line")


if __name__ == "__main__":
    app()
