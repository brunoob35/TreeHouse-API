<?php
$numeroCpf = $_POST['cpf'];

TransformarCpfEmArray($numeroCpf);

function ValidarCpf($array) {
    if (CheckIfAllEqual($array)) {
        if (CalcularDv1($array)) {
            if (CalcularDv2($array)) {
                $resultado = 'CPF VALIDO!';
                echo $resultado;
                return true;
            }
        }
    }
    $resultado = 'CPF INVALIDO!';
    echo $resultado;
    return false;
}

function CheckIfAllEqual($array) {
    $firstNumber = $array[0];
    for ($i = 1; $i < count($array); $i++) {
        if ($array[$i] !== $firstNumber) {
            return true;
        }
    }
    echo 'Numeros iguais <br>';
    return false;
}

function CalcularDv1($array) {
    $soma = 0;
    for ($i = 0; $i <= 8; $i++) {
        $multiplicacao = $array[$i] * (10 - $i);
        $soma += $multiplicacao;
    }
    $resultado = ($soma * 10) / 11;
    $resultadoArredonadoAcima = ceil($resultado * 10) / 10;

    $stringResultado = explode('.', strval($resultadoArredonadoAcima))[1];
    if ($stringResultado == $array[9]) {
        echo 'dv1 deu boa! <br>';
        return true;
    } else {
        return false;
    }
}

function CalcularDv2($array) {
    $soma = 0;
    for ($i = 0; $i <= 9; $i++) {
        $multiplicacao = $array[$i] * (11 - $i);
        $soma += $multiplicacao;
    }
    $resultado = ($soma * 10) / 11;
    $resultadoArredonadoAcima = ceil($resultado * 10) / 10;
    $stringResultado = explode('.', strval($resultadoArredonadoAcima))[1];
    if ($stringResultado == $array[10]) {
        echo 'dv2 deu boa! <br>';
        return true;
    } else {
        return false;
    }
}

function TransformarCpfEmArray($n) {
    $stringCpf = preg_replace('/\D/', "", $n);

    $arrayCpf = str_split($stringCpf);
    echo json_encode($arrayCpf).'<br>';
    ValidarCpf($arrayCpf);
}