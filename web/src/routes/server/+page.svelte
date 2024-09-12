<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
</script>

<p>Server Page (W.I.P)</p>

<p>Version: {PUBLIC_VERSION}</p>
<p>Commit: {PUBLIC_COMMIT}</p>

<div class="flex flex-col">
  <a href="/server/import">Import Album</a>
  <a href="/server/edit/album">Edit Albums</a>
</div>

<form action="?/runExport" method="post">
  <button>Run Export</button>
</form>
