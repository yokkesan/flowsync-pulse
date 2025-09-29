<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\Branch;
use App\Models\BranchLog;

class BranchController extends Controller
{
    public function index()
    {
        return Branch::latest()->get();
    }

    public function store(Request $request)
    {
        $validated = $request->validate([
            'project_name'  => 'required|string|max:100',   // user → project_name
            'branch_name'   => 'required|string|max:255',
            'ticket_number' => 'nullable|string|max:50',
        ]);

        // 1. 最新状態を branches に保存（プロジェクト単位で更新）
        $branch = Branch::updateOrCreate(
            ['project_name' => $validated['project_name']], // プロジェクトごとに1件
            [
                'branch_name'   => $validated['branch_name'],
                'ticket_number' => $validated['ticket_number'],
            ]
        );

        // 2. 履歴を branch_logs に追加
        BranchLog::create($validated);

        return response()->json([
            'message' => 'Saved successfully',
            'data'    => $branch,
        ], 201);
    }
}
