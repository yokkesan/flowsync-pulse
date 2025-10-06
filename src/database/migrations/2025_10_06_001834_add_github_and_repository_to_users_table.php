<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::table('users', function (Blueprint $table) {
            // GitHub アカウントとリポジトリ情報を追加
            $table->string('github')->nullable()->after('password');
            $table->string('repository')->nullable()->after('github');
        });
    }

    public function down(): void
    {
        Schema::table('users', function (Blueprint $table) {
            $table->dropColumn(['github', 'repository']);
        });
    }
};