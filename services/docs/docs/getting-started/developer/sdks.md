# Building SDK Packages 

We use [openapi-generator](https://openapi-generator.tech/config/generators) to create SDK packages using a spec taken directly from our C++ API

You can find the sdk builder, the settings, and the spec in the `sdk/` folder.

```shell
  --- SDK Generation ---
  make sdk                          Build all SDK packages
  make sdk android                  Build the android SDK package
  make sdk bash                     Build the bash SDK package
  make sdk cpp-oatpp-client         Build the cpp-oatpp-client SDK package
  make sdk cpp-oatpp-server         Build the cpp-oatpp-server SDK package
  make sdk dart                     Build the dart SDK package
  make sdk gdscript                 Build the gdscript SDK package
  make sdk go                       Build the go SDK package
  make sdk graphql-schema           Build the graphql-schema SDK package
  make sdk haskell-http-client      Build the haskell-http-client SDK package
  make sdk java                     Build the java SDK package
  make sdk jetbrains-http-client    Build the jetbrains-http-client SDK package
  make sdk k6                       Build the k6 SDK package
  make sdk lua                      Build the lua SDK package
  make sdk markdown                 Build the markdown SDK package
  make sdk mysql-schema             Build the mysql-schema SDK package
  make sdk nim                      Build the nim SDK package
  make sdk php                      Build the php SDK package
  make sdk powershell               Build the powershell SDK package
  make sdk protobuf-schema          Build the protobuf-schema SDK package
  make sdk python                   Build the python SDK package
  make sdk r                        Build the r SDK package
  make sdk ruby                     Build the ruby SDK package
  make sdk rust                     Build the rust SDK package
  make sdk swift5                   Build the swift5 SDK package
  make sdk swift6                   Build the swift6 SDK package
  make sdk typescript               Build the typescript SDK package
  make sdk typescript-angular       Build the typescript-angular SDK package
  make sdk typescript-node          Build the typescript-node SDK package
  make sdk wsdl-schema              Build the wsdl-schema SDK package
  make sdk zapier                   Build the zapier SDK package
```

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$blocks = new OpenAPI\Client\Api\BlockApi();

$hash = 'hash_example'; 

try {
    $block = $blocks->getBlockByHash($hash);
    print_r($block);
    print_r($block->transation_details[0]->pub_key);
} catch (Exception $e) {
    echo 'Exception when calling BlockApi->getBlockByHash: ', $e->getMessage(), PHP_EOL;
}
```