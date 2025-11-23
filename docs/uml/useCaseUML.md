flowchart LR
    %% Definição dos Atores (usando círculos duplos para diferenciar)
    Colab((Colaborador))
    Gestor((Gestor de RH))

    %% Fronteira do Sistema (System Boundary)
    subgraph System ["TimeSheet PRO"]
        direction TB

        %% Pacote de Autenticação
        subgraph Auth ["Autenticação"]
            UC1([Criar Conta])
            UC2([Fazer Login])
        end

        %% Pacote de Organização
        subgraph Org ["Gestão da Organização"]
            UC4([Criar Organização])
            UC5([Gerenciar Cargos])
            UC6([Convidar/Remover Usuários])
        end

        %% Pacote de Jornada
        subgraph Journey ["Jornada de Trabalho"]
            UC7([Registrar Ponto])
            UC8([Visualizar Espelho])
        end

        %% Pacote Admin
        subgraph Admin ["Painel Administrativo"]
            UC10([Dashboard da Equipe])
            UC11([Relatórios de Horas])
            UC12([Aprovar Horas Extras])
        end
    end

    %% Relacionamentos do Colaborador
    Colab --> UC1
    Colab --> UC2
    Colab --> UC7
    Colab --> UC8

    %% Relacionamentos do Gestor
    Gestor --> UC1
    Gestor --> UC2
    Gestor --> UC4
    Gestor --> UC5
    Gestor --> UC6
    Gestor --> UC10
    Gestor --> UC11
    Gestor --> UC12

    %% Estilização para diferenciar (Opcional)
    classDef actorStyle fill:#f9f,stroke:#333,stroke-width:2px;
    class Colab,Gestor actorStyle;
