<?php
// Converts every *.csv in --source into a same-named *.json in --dest.
// One-time/repeatable conversion step — this script is not shipped to consumers,
// only the JSON output in catalogs/ is. Run again whenever a source CSV changes.

declare(strict_types=1);

$options = getopt('', ['source:', 'dest:']);
$sourceDir = $options['source'] ?? null;
$destDir = $options['dest'] ?? null;

if ($sourceDir === null || $destDir === null) {
    fwrite(STDERR, "Usage: php csv_to_json.php --source=<dir with .csv files> --dest=<dir for .json output>\n");
    exit(1);
}

if (!is_dir($destDir) && !mkdir($destDir, 0777, true) && !is_dir($destDir)) {
    fwrite(STDERR, "Could not create destination directory: {$destDir}\n");
    exit(1);
}

$csvFiles = glob(rtrim($sourceDir, '/\\') . '/*.csv');
if ($csvFiles === false || $csvFiles === []) {
    fwrite(STDERR, "No .csv files found in {$sourceDir}\n");
    exit(1);
}

sort($csvFiles);

foreach ($csvFiles as $csvPath) {
    $name = pathinfo($csvPath, PATHINFO_FILENAME);
    $handle = fopen($csvPath, 'r');
    if ($handle === false) {
        fwrite(STDERR, "Could not open {$csvPath}\n");
        continue;
    }

    $header = fgetcsv($handle);
    if ($header === false) {
        fclose($handle);
        continue;
    }

    $rows = [];
    while (($row = fgetcsv($handle)) !== false) {
        if ($row === [null] || $row === false) {
            continue;
        }
        $rows[] = array_combine($header, $row);
    }
    fclose($handle);

    $jsonPath = rtrim($destDir, '/\\') . '/' . $name . '.json';
    $json = json_encode($rows, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);
    file_put_contents($jsonPath, $json . "\n");

    echo "  {$name}.csv -> {$name}.json (" . count($rows) . " rows)\n";
}

echo "Done.\n";
