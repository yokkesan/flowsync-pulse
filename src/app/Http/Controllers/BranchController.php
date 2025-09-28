<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\Branch;

class BranchController extends Controller
{
    public function index()
    {
        return Branch::latest()->get();
    }

    public function store(Request $request)
    {
        $validated = $request->validate([
            'user'          => 'nullable|string|max:100',
            'branch_name'   => 'required|string|max:255',
            'ticket_number' => 'nullable|string|max:50',
        ]);

        $branch = Branch::create($validated);

        return response()->json([
            'message' => 'Saved successfully',
            'data'    => $branch,
        ], 201);
    }
}