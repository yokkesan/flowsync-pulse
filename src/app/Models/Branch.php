<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Branch extends Model
{
    protected $fillable = [
    'user',
    'branch_name',
    'ticket_number',
    ];
}
